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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/huandu/go-clone"
	"github.com/remeh/sizedwaitgroup"
	"github.com/spf13/cobra"
)

var versionInfo string = "gfuzz v1.1"

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
		id := 0
		filterRequestsNum := 0
		errorRequestsNum := 0
		// ? calculate time
		startTime := time.Now()
		// ? show version
		// ? get params
		url, _ := cmd.Flags().GetString("url")
		timeout, _ := cmd.Flags().GetInt("timeout")
		mode, _ := cmd.Flags().GetString("mode")
		threadsNum, _ := cmd.Flags().GetInt("thread")
		totalPayloads, _ := cmd.Flags().GetStringArray("payload")
		totalHeaders, _ := cmd.Flags().GetStringArray("header")
		totalCookies, _ := cmd.Flags().GetStringArray("cookie")
		totalData, _ := cmd.Flags().GetStringArray("data")
		auth, _ := cmd.Flags().GetString("auth")
		showVerbose, _ := cmd.Flags().GetBool("verbose")
		requestsMethod, _ := cmd.Flags().GetString("method")

		// ? get filters
		filterAliasStrings := []string{"sc", "sh", "sw", "sl", "hc", "hh", "hw", "hl"}
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
		if strings.Contains(expression, "FUZZ") {
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
			// ! Stop if payloads set (number) invalid
			if lenOfstringList < 2 {
				utils.PrintError("Payloads set invalid")
				return
			}
			var payload payloads.Payload
			// payload := payloads.PayloadsBase{}
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
		tch := safechannels.New()
		tch.SafeSend(true)

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
		if !showVerbose {
			fmt.Println(`
===================================================================
ID           Response   Lines    Word     Chars       Payload
===================================================================`)
		} else {
			fmt.Println(`
===================================================================================================================
ID           C.Time       Response   Lines    Word     Chars       Payload      Md5Hash
===================================================================================================================`)
		}

		for range tch.Channel() {
			// * wait for total task
			tch.SafeSend(true)

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
								tch.SafeClose()
								wfp.SafeClose()
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
								tch.SafeClose()
								wfp.SafeClose()
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
								tch.SafeClose()
								wfp.SafeClose()
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
			// * continue if recv continueSignal
			select {
			case <-continueSignal.Channel():
				continue
			default:
			}
			// ! both channel are closed
			if wfp.IsClosed() && tch.IsClosed() {
				break
			}
			payloadDataCopy := clone.Clone(payloadsData).([]interface{})
			// * wait for sizedwaitgroup
			wg.Add()
			// plus id
			id++
			// ? replace placeholders to payload
			substring := "FUZZ"
			// clone from headers and cookies
			finalUrl := url
			finalAuth := auth
			finalExpression := expression
			finalData := clone.Clone(totalData).([]string)
			finalHeaders := clone.Clone(totalHeaders).([]string)
			finalCookies := clone.Clone(totalCookies).([]string)

			for index, payload := range payloadsData {
				if index > 0 {
					substring = "FUZ" + strconv.Itoa(index+1) + "Z"
				}
				// for url
				if strings.Contains(url, substring) {
					utils.ReplacePlaceHolderToPayload(index, &finalUrl, payload.(string), false)
				}

				// for auth
				if strings.Contains(auth, substring) {
					utils.ReplacePlaceHolderToPayload(index, &finalAuth, payload.(string), false)
				}

				// for expression
				if strings.Contains(expression, substring) {
					utils.ReplacePlaceHolderToPayload(index, &finalExpression, payload.(string), true)
				}

				// for data
				for i, s := range totalData {
					if strings.Contains(s, substring) {
						utils.ReplacePlaceHolderToPayloadFromArray(index, finalData, i, payload.(string))
					}
				}

				// for hedaers
				for i, s := range totalHeaders {
					if strings.Contains(s, substring) {
						utils.ReplacePlaceHolderToPayloadFromArray(index, finalHeaders, i, payload.(string))
					}
				}
				// for cookies
				for i, s := range totalCookies {
					if strings.Contains(s, substring) {
						utils.ReplacePlaceHolderToPayloadFromArray(index, finalCookies, i, payload.(string))
					}
				}

			}
			// ? test
			//fmt.Printf("Test %#v %#v %#v %#v %#v %#v\n", url, finalUrl, totalData, finalData, finalCookies, finalHeaders)

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
				resp, err := requests.Requests(requestsMethod, finalUrl, map[string][]string{"Headers": finalHeaders, "Cookies": finalCookies, "Data": finalData, "Auth": auth})
				if err != nil {
					utils.PrintErrorWithoutBlank("requests " + finalUrl + " error")
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
					if !showVerbose {
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
		fmt.Println()
		fmt.Printf("%s %.3fs\n", utils.Pcyan("Total time:         "), elapsed.Seconds())
		fmt.Println(utils.Pcyan("Processed Requests: "), id)
		fmt.Println(utils.Pcyan("Filtered Requests:  "), filterRequestsNum)
		fmt.Println(utils.Pcyan("Error Requests:     "), errorRequestsNum)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gfuzz.yaml)")
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
	rootCmd.Flags().StringArrayP("header", "H", []string{}, `Use header (ex:"Cookie:id=1312321&user=FUZZ"). Repeat option for various headers.`)
	rootCmd.Flags().StringArrayP("cookie", "b", []string{}, "Specify a cookie for the requests. Repeat option for various cookies.")
	rootCmd.Flags().String("auth", "", `in format "user:pass" or "FUZZ:FUZ2Z"`)
	rootCmd.Flags().StringP("method", "X", "GET", "Specify an HTTP method for the request, ie. HEAD or FUZZ, Support GET/POST/HEAD/OPTIONS/PUT/DELETE.")
	rootCmd.Flags().StringP("mode", "m", "zip", "Specify an iterator(zip, chain, product) for combining payloads.")
	rootCmd.Flags().StringArrayP("payload", "z", []string{}, "Specify a payload for each FUZZ keyword used in the form of name[,parameter][,encoder].")
	rootCmd.Flags().StringP("wordlist", "w", "", "The same as -z file, .")
	rootCmd.Flags().StringP("range", "r", "", "The same as -z range, .")
	rootCmd.Flags().String("list", "", "The same as -z list, .")
	rootCmd.Flags().String("filter", "", "Show/hide responses using the specified filter expression.")
	rootCmd.Flags().String("sc", "", "Show responses with the specified code.")
	rootCmd.Flags().String("hc", "", "Hide responses with the specified code.")
	rootCmd.Flags().String("sh", "", "Show responses with the specified chars.")
	rootCmd.Flags().String("hh", "", "Hide responses with the specified chars.")
	rootCmd.Flags().String("sw", "", "Show responses with the specified words.")
	rootCmd.Flags().String("hw", "", "Hide responses with the specified words.")
	rootCmd.Flags().String("sl", "", "Show responses with the specified lines.")
	rootCmd.Flags().String("hl", "", "Hide responses with the specified lines.")
	rootCmd.Flags().Bool("verbose", false, "Show verbose of fuzz results.")
	rootCmd.Flags().IntP("thread", "t", 16, "Threads of fuzz.")
	rootCmd.Flags().Int("timeout", 300, "Timeout second for fuzz.")
}
