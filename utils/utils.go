package utils

import (
	"fmt"
	"gfuzz/encoders"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var Pyellow = color.New(color.FgYellow).SprintFunc()
var Pred = color.New(color.FgRed).SprintFunc()
var Pcyan = color.New(color.FgCyan).SprintFunc()
var Pgreen = color.New(color.FgGreen).SprintFunc()

func PrintError(s string) {
	fmt.Printf("\n[%s] %s\n", Pred("-"), Pred(s))
}

func PrintWarn(s string) {
	fmt.Printf("\n[%s] %s\n", Pyellow("!"), Pyellow(s))
}

func PrintInfo(s string) {
	fmt.Printf("\n[%s] %s\n", Pcyan("#"), Pcyan(s))
}

func PrintSuccress(s string) {
	fmt.Printf("\n[%s] %s\n", Pgreen("+"), Pgreen(s))
}

func PrintErrorWithoutBlank(s string) {
	fmt.Printf("[%s] %s\n", Pred("-"), Pred(s))
}

func PrintRequestsError(id int, err error) {
	fmt.Printf("%s: %s\n", Pred(fmt.Sprintf("%09d", id)), Pred("Error "+err.Error()))
}

func PrintResponse(id int, responseCode int, linesNum int, wordsNum int, CharsNum int, payloads ...interface{}) {
	code := "000"
	if responseCode < 200 {
		code = Pcyan(strconv.Itoa(responseCode))
	} else if responseCode < 300 {
		code = Pgreen(strconv.Itoa(responseCode))
	} else if responseCode < 400 {
		code = Pyellow(strconv.Itoa(responseCode))
	} else {
		code = Pred(strconv.Itoa(responseCode))
	}
	fmt.Printf("%09d:   %-17s   %-3d L    %-3d W    %-3d Ch      \" ", id, code, linesNum, wordsNum, CharsNum)
	for _, payload := range payloads {
		fmt.Printf("%v ", payload)
	}
	fmt.Printf("\"\n")
}

func PrintResponseVerbose(id int, ctime time.Duration, responseCode int, linesNum int, wordsNum int, CharsNum int, md5hash string, payloads ...interface{}) {
	code := "000"
	if responseCode < 200 {
		code = Pcyan(strconv.Itoa(responseCode))
	} else if responseCode < 300 {
		code = Pgreen(strconv.Itoa(responseCode))
	} else if responseCode < 400 {
		code = Pyellow(strconv.Itoa(responseCode))
	} else {
		code = Pred(strconv.Itoa(responseCode))
	}
	fmt.Printf("%09d:   %-.3fs       %-17s   %-3d L    %-3d W    %-3d Ch      \" ", id, ctime.Seconds(), code, linesNum, wordsNum, CharsNum)
	for _, payload := range payloads {
		fmt.Printf("%v ", payload)
	}
	fmt.Printf(`"        ` + md5hash + "\n")
}

func CalcplaceHoldersNum(num *int, s string) {
	regc, _ := regexp.Compile("FUZ(\\d*)Z")
	matches := regc.FindAllStringSubmatch(s, -1)
	lenOfMatches := len(matches)
	if lenOfMatches < 1 {
		return
	}
	placeHolderString := matches[len(matches)-1][1]
	tempNum, err := strconv.Atoi(placeHolderString)
	if err != nil {
		tempNum = 1
	}
	if tempNum > *num {
		*num = tempNum
	}
}

func ReplacePlaceHolderToPayload(id int, s *string, payload string, isExpression bool) {
	if isExpression {
		_, err := strconv.Atoi(payload)
		if err != nil {
			payload = "'" + payload + "'"
		}
	}
	if id == 0 {
		*s = strings.ReplaceAll(*s, "FUZZ", payload)
	} else {
		*s = strings.ReplaceAll(*s, "FUZ"+strconv.Itoa(id+1)+"Z", payload)
	}

}

func ReplacePlaceHolderToPayloadFromArray(id int, stringArray []string, arrayIndex int, payload string) {
	if len(stringArray) < 0 {
		return
	}
	if id == 0 {
		stringArray[arrayIndex] = strings.ReplaceAll(stringArray[arrayIndex], "FUZZ", payload)
	} else {
		stringArray[arrayIndex] = strings.ReplaceAll(stringArray[arrayIndex], "FUZ"+strconv.Itoa(id+1)+"Z", payload)
	}

}

func ProductStringWithStrings(sets ...[]interface{}) chan interface{} {
	ch := make(chan interface{}, 0)
	lens := func(i int) int { return len(sets[i]) }
	go func(lens func(i int) int, sets ...[]interface{}) {
		for ix := make([]int, len(sets)); ix[0] < lens(0); nextIndex(ix, lens) {
			for j, k := range ix {
				ch <- sets[j][k].(string)
			}
		}
		close(ch)
	}(lens, sets...)
	return ch
}

// ? product string with runes by sets
func ProductStringWithRunes(sets ...[]interface{}) chan interface{} {
	ch := make(chan interface{}, 0)
	lens := func(i int) int { return len(sets[i]) }
	go func(lens func(i int) int, sets ...[]interface{}) {
		for ix := make([]int, len(sets)); ix[0] < lens(0); nextIndex(ix, lens) {
			var r []interface{}
			for j, k := range ix {
				r = append(r, sets[j][k])
			}
			res := ""
			for _, ch := range r {
				res += ch.(string)
			}
			ch <- res
		}
		close(ch)
	}(lens, sets...)

	return ch
}

func nextIndex(ix []int, lens func(i int) int) {
	for j := len(ix) - 1; j >= 0; j-- {
		ix[j]++
		if j == 0 || ix[j] < lens(j) {
			return
		}
		ix[j] = 0
	}
}

// ? encode payloads by call encoders
func EncodeForPayload(encoderArray []string, payload interface{}) interface{} {
	if len(encoderArray) > 0 {
		for _, e := range encoderArray {
			payload = encoders.GetEncoder(e).Encode(payload)
		}
	}
	return payload
}
