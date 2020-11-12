package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	re, _ := regexp.Compile("\\[\\d+m")
	for i, s := range a {
		str, ok := s.(string)
		if ok {
			a[i] = re.ReplaceAllString(str, "")
		}
	}
	return a
}

func Println(a ...interface{}) {
	fmt.Println(a...)
	if logFile != nil {
		fmt.Fprintln(logFile, cleanColor(a)...)
	}

}

func Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	if logFile != nil {
		fmt.Fprintf(logFile, format, cleanColor(a)...)
	}
}

func Fprintf(format string, a ...interface{}) {
	if logFile != nil {
		fmt.Fprintf(logFile, format, cleanColor(a)...)
	}
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
	fileFormat := `%09d:   %-8s   %-3d L    %-3d W    %-3d Ch      %-s`
	format := `%09d:   %-17s   %-3d L    %-3d W    %-3d Ch      %-s`
	payloadsString := `" `
	for _, payload := range payloads {
		payloadsString += fmt.Sprintf("%v ", payload)
	}
	payloadsString += `"`
	a := []interface{}{id, code, linesNum, wordsNum, CharsNum, payloadsString}
	fa := []interface{}{id, strconv.Itoa(responseCode), linesNum, wordsNum, CharsNum, payloadsString}

	fmt.Printf(format, a...)
	fmt.Println()
	Fprintf(fileFormat, fa...)
	Fprintf("\n")
}

func PrintResponseVerbose(id int, ctime time.Duration, responseCode int, linesNum int, wordsNum int, CharsNum int, md5hash string, payloads ...interface{}) {
	printLock.Lock()
	defer printLock.Unlock()
	code := ColorStatusCode(responseCode)
	fileFormat := `%09d:   %-.3fs       %-8s   %-3d L    %-3d W    %-3d Ch      %-30s         %s`
	format := `%09d:   %-.3fs       %-17s   %-3d L    %-3d W    %-3d Ch      %-30s         %s`
	payloadsString := `" `
	for _, payload := range payloads {
		payloadsString += fmt.Sprintf("%v ", payload)
	}
	payloadsString += `"`
	a := []interface{}{id, ctime.Seconds(), code, linesNum, wordsNum, CharsNum, payloadsString, md5hash}
	fa := []interface{}{id, ctime.Seconds(), strconv.Itoa(responseCode), linesNum, wordsNum, CharsNum, payloadsString, md5hash}

	fmt.Printf(format, a...)
	fmt.Println()
	Fprintf(fileFormat, fa...)
	Fprintf("\n")
}

func PrintTips(isShowVerbose bool) {
	if !isShowVerbose {
		format := `%-10s   %-8s   %-5s    %-6s    %-6s     %-s`
		tip := fmt.Sprintf(format, "ID", "Response", "Lines", "Words", "Chars", "Payload")
		line := strings.Repeat("=", len(tip)+10)
		Println(line + "\n" + tip + "\n" + line)
	} else {
		format := `%-10s   %-5s       %-8s   %-5s    %-6s    %-6s     %-30s         %s`
		tip := fmt.Sprintf(format, "ID", "C.Time", "Response", "Lines", "Words", "Chars", "Payload", "Md5Hash")
		line := strings.Repeat("=", len(tip)+28)
		Println(line + "\n" + tip + "\n" + line)
	}
}
