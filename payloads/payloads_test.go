package payloads

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestRange(t *testing.T) {
	// ? test 0-9
	a := PayloadRange{}
	a.New("0-9")
	ch := a.Channel()
	sets := make([]string, 10, 10)
	i := 0
	for set := range ch {
		sets[i] = set.(string)
		i++
	}
	for i := 0; i <= 9; i++ {
		if sets[i] != strconv.Itoa(i) {
			t.Fatalf("range numbers expected = %d, actual: %s", i, sets[i])
		}
	}

	// ? tst aa-dd
	b := PayloadRange{}
	b.New("aa-dd")
	ch = b.Channel()
	sets2 := make([]interface{}, 4*4)
	i = 0
	for set := range ch {
		sets2[i] = set.(string)
		i++
	}
	i = 0
	for y := rune('a'); y <= rune('d'); y++ {
		for x := rune('a'); x <= rune('d'); x++ {
			actual := ""
			expected := string(y) + string(x)
			actual = sets2[i].(string)
			i++
			if actual != expected {
				t.Fatalf("range chars expected = %c%c, actual: %s", y, x, actual)
			}
		}
	}
}

func TestFile(t *testing.T) {
	// ? test file_test.txt
	a := PayloadFile{}
	a.New("./file_test.txt")
	ch := a.Channel()
	i := 0
	bytes, err := ioutil.ReadFile("./file_test.txt")
	if err != nil {
		t.Error("file not exists")
	}
	str := string(bytes)
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.ReplaceAll(str, "\r", "")
	strs := strings.Split(str, "\n")
	for actual := range ch {
		if strs[i] != actual {
			t.Fatalf("file numbers expected = %s, actual: %s", strs[i], actual)
		}
		i++
	}
}
