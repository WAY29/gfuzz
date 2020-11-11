package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/color"
)

var Pyellow = color.New(color.FgYellow).SprintFunc()
var Pred = color.New(color.FgRed).SprintFunc()
var Pcyan = color.New(color.FgCyan).SprintFunc()
var Pgreen = color.New(color.FgGreen).SprintFunc()
var logFile *os.File
var printLock sync.Mutex

func SetWriter(f *os.File) {
	logFile = f
}

func cleanColor(a []interface{}) []interface{} {
	// a := make([]interface{}, len(a))
	re, _ := regexp.Compile("\\[\\d+m")
	for i, s := range a {
		// fmt.Printf("test %#v\n", s)
		str, ok := s.(string)
		if ok {
			a[i] = re.ReplaceAllString(str, "")
		}
	}
	return a
}

func Println(a ...interface{}) {
	printLock.Lock()
	defer printLock.Unlock()
	fmt.Println(a...)
	fmt.Fprintln(logFile, cleanColor(a)...)
}

func Printf(format string, a ...interface{}) {
	printLock.Lock()
	defer printLock.Unlock()
	// fmt.Printf("Test %s %#v\n", format, cleanColor(a))
	fmt.Printf(format, a...)
	fmt.Fprintf(logFile, format, cleanColor(a)...)
}

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
	printLock.Lock()
	defer printLock.Unlock()
	code := ColorStatusCode(responseCode)
	fileFormat := "%09d:   %-8s   %-3d L    %-3d W    %-3d Ch      \" "
	format := "%09d:   %-17s   %-3d L    %-3d W    %-3d Ch      \" "
	a := []interface{}{id, code, linesNum, wordsNum, CharsNum}
	fa := []interface{}{id, strconv.Itoa(responseCode), linesNum, wordsNum, CharsNum}
	fmt.Printf(format, a...)
	fmt.Fprintf(logFile, fileFormat, fa...)
	for _, payload := range payloads {
		fmt.Printf("%v ", payload)
		fmt.Fprintf(logFile, "%v ", payload)

	}
	fmt.Printf(" \"\n")
	fmt.Fprintf(logFile, " \"\n")
}

func PrintResponseVerbose(id int, ctime time.Duration, responseCode int, linesNum int, wordsNum int, CharsNum int, md5hash string, payloads ...interface{}) {
	printLock.Lock()
	defer printLock.Unlock()
	code := ColorStatusCode(responseCode)
	fileFormat := "%09d:   %-.3fs       %-8s   %-3d L    %-3d W    %-3d Ch      \" "
	format := "%09d:   %-.3fs       %-17s   %-3d L    %-3d W    %-3d Ch      \" "
	a := []interface{}{id, ctime.Seconds(), code, linesNum, wordsNum, CharsNum}
	fa := []interface{}{id, ctime.Seconds(), strconv.Itoa(responseCode), linesNum, wordsNum, CharsNum}
	fmt.Printf(format, a...)
	fmt.Fprintf(logFile, fileFormat, fa...)
	for _, payload := range payloads {
		fmt.Printf("%v ", payload)
		fmt.Fprintf(logFile, "%v ", payload)
	}
	fmt.Printf(` "        ` + md5hash + "\n")
	fmt.Fprintf(logFile, ` "        `+md5hash+"\n")
}
