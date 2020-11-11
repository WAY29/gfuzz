package payloads

import (
	"errors"
	"gfuzz/utils"
	"strconv"
	"strings"
)

func init() {
	AddPayload("range", &PayloadRange{})
	AddPayloadInfo("range", "Returns each number or char of the given range.")
}

type PayloadRange struct {
	channel chan interface{}
}

func isNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *PayloadRange) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadRange) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	tempArray := strings.Split(s[0].(string), "-")
	if len(tempArray) > 0 {
		if isNum(tempArray[0]) { // numbers
			startNum, ok1 := strconv.Atoi(tempArray[0])
			endNum, ok2 := strconv.Atoi(tempArray[1])
			if ok1 == nil && ok2 == nil { // both is valid

			} else {
				return errors.New("Invaild numbers")
			}
			if startNum > endNum { // Var swap
				startNum, endNum = endNum, startNum
			}
			go func() {
				defer func() {
					close(ch)
				}()
				for i := startNum; i <= endNum; i++ {
					ch <- strconv.Itoa(i)
				}
			}()
		} else { // chars
			if len(tempArray[0]) != len(tempArray[1]) { // both two string length must be same
				return errors.New("Invaild payloads length")
			}
			sets := make([][]interface{}, len(tempArray[0]))
			for i, c := range tempArray[0] {
				c1 := c
				c2 := rune(tempArray[1][i])
				if c1 > c2 {
					c1, c2 = c2, c1
				}
				for cc := c1; cc <= c2; cc++ {
					sets[i] = append(sets[i], string(cc))
				}
			}
			p.channel = utils.ProductStringWithRunes(sets...)
		}

	} else {
		return errors.New("Invaild payloads")
	}

	return nil
}
