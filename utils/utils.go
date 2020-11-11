package utils

import (
	"fmt"
	"gfuzz/encoders"
	"regexp"
	"strconv"
	"strings"
)

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

// ? python-like format functions
func Format(format string, p map[string]interface{}) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = k
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

// ? python-like format functions
func FormatStringArray(stringArray []string, p map[string]interface{}) []string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = k
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	for i, s := range stringArray {
		stringArray[i] = strings.NewReplacer(args...).Replace(s)
	}
	return stringArray
}

// ? color StatusCode
func ColorStatusCode(responseCode int) string {
	code := fmt.Sprintf("%03s", strconv.Itoa(responseCode))
	if 100 < responseCode && responseCode < 200 {
		code = Pcyan(code)
	} else if responseCode < 300 {
		code = Pgreen(code)
	} else if responseCode < 400 {
		code = Pyellow(code)
	} else {
		code = Pred(code)
	}
	return code
}
