/*
TODO features:
* y other requests methods
output file
* y verbose
* y auth

*/
package cmd

import (
	"fmt"
	"gfuzz/encoders"
	"gfuzz/filters"
	"gfuzz/payloads"
	"gfuzz/requests"
	"gfuzz/safechannels"
	"gfuzz/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/huandu/go-clone"
	"github.com/remeh/sizedwaitgroup"
	"github.com/spf13/cobra"
)

var versionInfo string = "gfuzz v1.5"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gfuzz",
	Short: "wfuzz copied to Golang",
	Long: `wfuzz copied to Golang.

Example: 
	gfuzz -u "http://httpbin.org/get?a=FUZZ&b=FUZZ" -m chain -z range,0-9 -r a-b # total 12 requests
	gfuzz -u "http://httpbin.org/get?a=FUZZ&b=FUZ2Z" -m zip -r 0-9 -r a-j # total 10 requests, zip mode is default
	gfuzz -u "http://httpbin.org/post?a=FUZZ" -d "b=FUZ2Z" -m product -r 0-9 -r a-b # total 10*2=20 requests`,
	Version: versionInfo,
	Run: func(cmd *cobra.Command, args []string) {
		// ? set vars
		id := 0
		filterRequestsNum := 0
		errorRequestsNum := 0
		var logFile *os.File
		// ? close logFile if set output
		defer func() {
			r := recover()
			if r != nil {
				utils.PrintError(r.(string))
			}
			if logFile != nil {
				logFile.Close()
			}
		}()
		// ? calculate time
		startTime := time.Now()
		// ? show version
		// ? get params
		url, _ := cmd.Flags().GetString("url")
		isUseSession, _ := cmd.Flags().GetBool("session")
		isFollow, _ := cmd.Flags().GetBool("follow")
		timeout, _ := cmd.Flags().GetInt("timeout")
		reqDelayTimeout, _ := cmd.Flags().GetInt("req_delay")
		connDelayTimeout, _ := cmd.Flags().GetInt("conn_delay")
		mode, _ := cmd.Flags().GetString("mode")
		outputFile, _ := cmd.Flags().GetString("output")
		threadsNum, _ := cmd.Flags().GetInt("thread")
		totalPayloads, _ := cmd.Flags().GetStringArray("payload")
		totalHeaders, _ := cmd.Flags().GetStringArray("header")
		totalCookies, _ := cmd.Flags().GetStringArray("cookie")
		totalData, _ := cmd.Flags().GetStringArray("data")
		totalJson, _ := cmd.Flags().GetString("json")
		auth, _ := cmd.Flags().GetString("auth")
		requestsMethod, _ := cmd.Flags().GetString("method")
		isShowVerbose, _ := cmd.Flags().GetBool("verbose")
		// ? read json data if startwith @
		if len(totalJson) > 0 && totalJson[0] == '@' {
			jsonDataBytes, err := ioutil.ReadFile(totalJson[1:])
			if err == nil {
				totalJson = strings.ReplaceAll(string(jsonDataBytes), "\n", "")
				totalJson = strings.ReplaceAll(totalJson, "\r", "")
			} else {
				totalJson = ""
			}
		}
		// ? set multi writers if output file
		if len(outputFile) > 0 {
			var err error
			logFile, err = os.OpenFile(outputFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}
			utils.SetWriter(logFile)
		}
		// ? get alias filters
		filterAliasStrings := []string{"sc", "sh", "sw", "sl", "hc", "hh", "hw", "hl", "sx", "hx"}
		filtersMap := make(map[string]string, len(filterAliasStrings))
		expression, _ := cmd.Flags().GetString("filter")
		for _, s := range filterAliasStrings {
			filtersMap[s], _ = cmd.Flags().GetString(s)
		}
		// ? add expression
		for k, v := range filtersMap {
			if len(v) > 0 {
				expression = filters.GetFilter(string(k[0])).AddExpression(expression, string(k[1]), v)
			}
		}

		// ! Stop if auth invaild
		if auth != "" && !strings.Contains(auth, ":") {
			utils.PrintError("Auth Invaild")
			return
		}
		// ! Stop if not contains placeholder
		check := false
		if strings.Contains(url, "FUZZ") {
			check = true
		}
		if strings.Contains(auth, "FUZZ") {
			check = true
		}
		if strings.Contains(requestsMethod, "FUZZ") {
			check = true
		}
		if strings.Contains(expression, "FUZZ") {
			check = true
		}
		if strings.Contains(totalJson, "FUZZ") {
			check = true
		}
		for _, s := range totalData {
			if strings.Contains(s, "FUZZ") {
				check = true
			}
		}
		for _, s := range totalHeaders {
			if strings.Contains(s, "FUZZ") {
				check = true
			}
		}
		for _, s := range totalCookies {
			if strings.Contains(s, "FUZZ") {
				check = true
			}
		}
		if !check {
			utils.PrintError("Must have FUZZ placeholder")
			return
		}
		// ? calc placeHoldersNum
		lenOfPayloads := len(totalPayloads)
		placeHoldersNum := 0

		// from url
		utils.CalcplaceHoldersNum(&placeHoldersNum, url)

		// from auth
		utils.CalcplaceHoldersNum(&placeHoldersNum, auth)

		// from requests method
		utils.CalcplaceHoldersNum(&placeHoldersNum, requestsMethod)

		// from filters expression
		utils.CalcplaceHoldersNum(&placeHoldersNum, expression)

		// from json
		utils.CalcplaceHoldersNum(&placeHoldersNum, totalJson)

		// from data
		for _, s := range totalData {
			utils.CalcplaceHoldersNum(&placeHoldersNum, s)
		}

		// from headers
		for _, s := range totalHeaders {
			utils.CalcplaceHoldersNum(&placeHoldersNum, s)
		}

		// from cookies
		for _, s := range totalCookies {
			utils.CalcplaceHoldersNum(&placeHoldersNum, s)
		}

		switch mode {
		// ! (zip/product mode) Stop if payloads number not same as placeholders number
		case "zip":
			fallthrough
		case "product":
			if lenOfPayloads != placeHoldersNum {
				utils.PrintError("[zip mode] Payloads number must be same as placeholders number")
				return
			}
		// ! (chain mode) Stop if placeholders number != 1
		case "chain":
			if placeHoldersNum != 1 {
				utils.PrintError("[chain mode] placeholders number must be 1")
				return
			}
		default:
			utils.PrintError("Unsupport mode")
			return
		}

		// ? parse payloads
		payloadsData := make([]interface{}, placeHoldersNum)
		channelsArray := make([]chan interface{}, lenOfPayloads)
		encodersArray := make([][]string, lenOfPayloads)

		for i, z := range totalPayloads {
			stringList := strings.Split(z, ",")
			lenOfstringList := len(stringList)
			// ? add empty string if len of stringlist less than 2
			if lenOfstringList < 2 {
				lenOfstringList = 2
				stringList = append(stringList, "")
			}
			var payload payloads.Payload

			// ? generate payloads
			payload = payloads.GetPayload(stringList[0])
			if payload == nil {
				utils.PrintError("Unsupport Payloads")
				return
			}

			err := payload.New(stringList[1])
			if err != nil {
				utils.PrintError("Payload generate Error")
				return
			}
			// * set payload channel
			channelsArray[i] = payload.Channel()
			// * set encoders
			if lenOfstringList > 2 {
				encodersArray[i] = stringList[2:]
			}

		}

		// channel index for chain mode
		chIndex := 0

		// prepare product vars for product mode
		productChannel := make(chan interface{}, 0)

		// prepare sizedwaitgroup
		wg := sizedwaitgroup.New(threadsNum)

		// wait for paylaodsdata channel
		wfp := safechannels.New()

		// wait for total task
		// tch := safechannels.New()
		stopFlag := false

		// continue signal
		continueSignal := safechannels.New()

		// ? prepare for product mode
		if mode == "product" {
			var payload interface{}
			ok := true
			utils.PrintInfo("Generate products...")
			tempPayloadsArrays := make([][]interface{}, lenOfPayloads)
			for i := 0; i < placeHoldersNum; i++ {
				exit := false
				tempPayloadsArray := make([]interface{}, 0)
				for !exit {
					select {
					case payload, ok = <-channelsArray[i]:
						if !ok || payload == nil { // read fromo channel error , maybe done
							exit = true
							break
						}
						tempPayloadsArray = append(tempPayloadsArray, payload)
					// ! stop if wait payloadsdata timeout
					case <-time.After(time.Duration(timeout) * time.Second):
						utils.PrintError("Read data timeout")
						return
					}
				}
				tempPayloadsArrays[i] = tempPayloadsArray

			}
			utils.PrintSuccress("Generate products finish")
			productChannel = utils.ProductStringWithStrings(tempPayloadsArrays...)
		}

		// ? print fuzz tips
		utils.PrintTips(isShowVerbose)

		for {
			// ? get paylaods data from channels
			go func(payloadsData []interface{}) {
				var payload interface{}
				ok := true
				switch mode {
				// ? zip mode
				case "zip":
					for i := 0; i < placeHoldersNum; i++ {
						select {
						case payload, ok = <-channelsArray[i]:

							if !ok || payload == nil { // read fromo channel error , maybe done
								// ! stop if payload is null
								wfp.SafeSend(true)
								stopFlag = true
								return
							}
							// ? generate encoders
							// * call encoders
							payload = utils.EncodeForPayload(encodersArray[i], payload)
							// * set payloads data
							payloadsData[i] = payload
						// ! stop if wait payloadsdata timeout
						case <-time.After(time.Duration(timeout) * time.Second):
							utils.PrintError("Read data timeout")
							return
						}
					}
					wfp.SafeSend(true)
				// ? chain mode
				case "chain":
					ch := channelsArray[chIndex]
					select {
					case payload, ok = <-ch:
						if !ok || payload == nil { // ! stop if all channel read
							chIndex++
							if chIndex >= lenOfPayloads {
								wfp.SafeSend(true)
								stopFlag = true
								return
							} else {
								continueSignal.SafeSend(true)
							}
						}
						// ? generate encoders
						// * call encoders
						payload = utils.EncodeForPayload(encodersArray[chIndex], payload)
						// * set payloads data
						payloadsData[0] = payload
						// ! stop if wait payloadsdata timeout
					case <-time.After(time.Duration(timeout) * time.Second):
						utils.PrintError("Read data timeout")
						return
					}
					wfp.SafeSend(true)
				//? product mode
				case "product":
					// ? fill payloads data from product
					for i := 0; i < placeHoldersNum; i++ {
						select {
						case payload, ok = <-productChannel:

							if !ok || payload == nil { // read fromo channel error , maybe done
								wfp.SafeSend(true)
								stopFlag = true
								return
							}
							// ? generate encoders
							// * call encoders
							payload = utils.EncodeForPayload(encodersArray[i], payload)
							// * set payloads data
							payloadsData[i] = payload
						// ! stop if wait payloadsdata timeout
						case <-time.After(time.Duration(timeout) * time.Second):
							utils.PrintError("Read data timeout")
							return
						}
					}
					wfp.SafeSend(true)
				}

			}(payloadsData)
			// * wait for payloads data
			select {
			case <-wfp.Channel():
			// ! stop if wait payloadsdata timeout
			case <-time.After(time.Duration(timeout) * time.Second):
				utils.PrintError("Read data timeout")
				return
			}
			// ! stop if no payloads
			if stopFlag {
				break
			}
			// * continue if recv continueSignal
			select {
			case <-continueSignal.Channel():
				continue
			default:
			}

			// if wfp.IsClosed() && tch.IsClosed() {
			// 	break
			// }
			payloadDataCopy := clone.Clone(payloadsData).([]interface{})
			// * wait for sizedwaitgroup
			wg.Add()
			// plus id
			id++
			// ? replace placeholders to payload
			placeholder := "FUZZ"
			finalPayloadsData := make(map[string]interface{}, lenOfPayloads*2)
			for index, payload := range payloadsData {
				if index > 0 {
					placeholder = "FUZ" + strconv.Itoa(index+1) + "Z"
				}
				finalPayloadsData[placeholder] = payload
			}
			finalUrl := utils.Format(url, finalPayloadsData)
			finalAuth := utils.Format(auth, finalPayloadsData)
			finalRequestMethod := utils.Format(requestsMethod, finalPayloadsData)
			finalJson := utils.Format(totalJson, finalPayloadsData)
			finalExpression := utils.Format(expression, finalPayloadsData)
			finalData := utils.FormatStringArray(clone.Clone(totalData).([]string), finalPayloadsData)
			finalHeaders := utils.FormatStringArray(clone.Clone(totalHeaders).([]string), finalPayloadsData)
			finalCookies := utils.FormatStringArray(clone.Clone(totalCookies).([]string), finalPayloadsData)

			// ? test
			// fmt.Printf("Test %#v %#v %#v %#v %#v %#v\n", url, finalUrl, totalData, finalData, finalCookies, finalHeaders)

			// ? start to fuzz
			go func(id int, payloadsData []interface{}) {
				startRequestTime := time.Now()
				defer func() {
					wg.Done()
				}()
				auth := []string{}
				if strings.Contains(finalAuth, ":") {
					auth = strings.Split(finalAuth, ":")
				}
				resp, err := requests.Requests(finalRequestMethod, finalUrl, map[string][]string{"Headers": finalHeaders, "Cookies": finalCookies, "Data": finalData, "Auth": auth},
					map[string]interface{}{"UseSession": isUseSession, "ReqTimeout": reqDelayTimeout, "ConnTimeout": connDelayTimeout, "Follow": isFollow, "Json": finalJson})
				if err != nil {
					if isShowVerbose {
						utils.PrintErrorWithoutBlank("requests " + finalUrl + " error " + err.Error())
					} else {
						utils.PrintErrorWithoutBlank("requests " + finalUrl + " error")
					}
					errorRequestsNum++
					return
				}
				// ?get response data
				text := resp.String()
				md5HashOftext := encoders.GetEncoder("md5").Encode(text).(string)
				lenOftext := len(text)
				lines := strings.Split(text, "\n")
				lenOflines := len(lines)
				words := strings.Split(text, " ")
				lenOfwords := len(words)
				isPass := true
				// ? filter response
				if expression != "" {
					parameters := make(map[string]interface{}, 32)
					parameters["c"] = resp.StatusCode
					parameters["code"] = resp.StatusCode
					parameters["h"] = lenOftext
					parameters["chars"] = lenOftext
					parameters["l"] = lenOflines
					parameters["lines"] = lenOflines
					parameters["w"] = lenOfwords
					parameters["words"] = lenOfwords
					parameters["id"] = id
					parameters["md5"] = md5HashOftext
					// ? --------------------------------
					parameters["url"] = resp.RawResponse.Request.URL.String()
					parameters["method"] = resp.RawResponse.Request.Method
					parameters["scheme"] = resp.RawResponse.Request.URL.Scheme
					parameters["host"] = resp.RawResponse.Request.URL.Host
					parameters["x"] = text
					parameters["content"] = text
					for _, cookie := range resp.RawResponse.Request.Cookies() {
						parameters["req_cookies_"+cookie.Name] = cookie.Value
					}
					for _, cookie := range resp.RawResponse.Cookies() {
						parameters["res_cookies_"+cookie.Name] = cookie.Value
					}
					for k, v := range resp.RawResponse.Request.Header {
						parameters["req_headers_"+k] = strings.Join(v, " ")
					}
					for k, v := range resp.RawResponse.Header {
						parameters["res_headers_"+k] = strings.Join(v, " ")
					}
					// ! test
					// fmt.Println("test", finalExpression)
					isPass, err = filters.FinalFilter(finalExpression, parameters)
					if err != nil {
						errorRequestsNum++
						utils.PrintRequestsError(id, err)
						return
					}
				}
				elapsedOfRequest := time.Since(startRequestTime)
				if isPass {
					if !isShowVerbose {
						utils.PrintResponse(id, resp.StatusCode, lenOflines, lenOfwords, lenOftext, payloadsData...)
					} else {
						utils.PrintResponseVerbose(id, elapsedOfRequest, resp.StatusCode, lenOflines, lenOfwords, lenOftext, md5HashOftext, payloadsData...)
					}
				} else {
					filterRequestsNum++
				}
			}(id, payloadDataCopy)

		}
		wg.Wait()
		elapsed := time.Since(startTime)
		utils.Println()
		utils.Printf("%s %.3fs\n", utils.Pcyan("Total time:         "), elapsed.Seconds())
		utils.Println(utils.Pcyan("Processed Requests: "), id)
		utils.Println(utils.Pcyan("Filtered Requests:  "), filterRequestsNum)
		utils.Println(utils.Pcyan("Error Requests:     "), errorRequestsNum)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// ? replace some alias flags to real flags
	args := os.Args[1:]
	lenOfArgs := len(args)
	for i := 0; i < lenOfArgs; i++ {
		if i+1 < lenOfArgs {
			if args[i] == "-w" || args[i] == "--wordlist" {
				args[i] = "-z"
				args[i+1] = "file," + args[i+1]
			} else if args[i] == "-r" || args[i] == "--range" {
				args[i] = "-z"
				args[i+1] = "range," + args[i+1]
			} else if args[i] == "--list" {
				args[i] = "-z"
				args[i+1] = "list," + args[i+1]
			}
		}
	}
	rootCmd.SetArgs(args)
	rootCmd.Flags().SortFlags = false
	// ? set flags
	rootCmd.Flags().StringP("url", "u", "", "(required) Target URL.")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().StringArrayP("data", "d", []string{}, `Use post data (ex: "id=FUZZ&catalogue=1"). Repeat option for various data.`)
	rootCmd.Flags().String("json", "", `Use json data post. Use @filepath to read json data from file`)
	rootCmd.Flags().StringArrayP("header", "H", []string{}, `Use header (ex:"Cookie:id=1312321&user=FUZZ"). Repeat option for various headers.`)
	rootCmd.Flags().StringArrayP("cookie", "b", []string{}, "Specify a cookie for the requests. Repeat option for various cookies.")
	rootCmd.Flags().String("auth", "", `in format "user:pass" or "FUZZ:FUZ2Z"`)
	rootCmd.Flags().StringP("method", "X", "GET", "Specify an HTTP method for the request, ie. HEAD or FUZZ, Support GET/POST/HEAD/OPTIONS/PUT/DELETE.")
	rootCmd.Flags().StringP("mode", "m", "zip", "Specify an iterator(zip, chain, product) for combining payloads.")
	rootCmd.Flags().StringArrayP("payload", "z", []string{}, "Specify a payload for each FUZZ keyword used in the form of name[,parameter][,encoder].")
	rootCmd.Flags().StringP("wordlist", "w", "", "The same as -z file, .")
	rootCmd.Flags().StringP("range", "r", "", "The same as -z range, .")
	rootCmd.Flags().String("list", "", "The same as -z list, .")
	rootCmd.Flags().StringP("filter", "f", "", "Show/hide responses using the specified filter expression.")
	rootCmd.Flags().String("sc", "", "Show responses with the specified code.")
	rootCmd.Flags().String("hc", "", "Hide responses with the specified code.")
	rootCmd.Flags().String("sh", "", "Show responses with the specified chars.")
	rootCmd.Flags().String("hh", "", "Hide responses with the specified chars.")
	rootCmd.Flags().String("sw", "", "Show responses with the specified words.")
	rootCmd.Flags().String("hw", "", "Hide responses with the specified words.")
	rootCmd.Flags().String("sl", "", "Show responses with the specified lines.")
	rootCmd.Flags().String("hl", "", "Hide responses with the specified lines.")
	rootCmd.Flags().String("sx", "", "Show responses with the specified content.")
	rootCmd.Flags().String("hx", "", "Hide responses with the specified content.")
	rootCmd.Flags().StringP("output", "o", "", "Store results in the output file using the specified printer")
	rootCmd.Flags().Bool("verbose", false, "Show verbose of fuzz results.")
	rootCmd.Flags().BoolP("session", "S", false, "Whether use session for fuzz.")
	rootCmd.Flags().BoolP("follow", "L", false, "Follow HTTP redirections.")
	rootCmd.Flags().IntP("thread", "t", 32, "Threads of fuzz.")
	rootCmd.Flags().Int("timeout", 300, "Timeout second for fuzz.")
	rootCmd.Flags().Int("req_delay", 90, "Sets the maximum time in seconds the request is allowed to take (CURLOPT_TIMEOUT).")
	rootCmd.Flags().Int("conn_delay", 90, "Sets the maximum time in seconds the connection phase to the server to take (CURLOPT_CONNECTTIMEOUT).")
}
