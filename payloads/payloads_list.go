package payloads

import (
	"strings"
)

func init() {
	AddPayload("list", &PayloadList{})
	AddPayloadInfo("list", "Returns each element of the given word list separated by -.")
}

type PayloadList struct {
	channel chan interface{}
}

//
func (p *PayloadList) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadList) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	str := s[0].(string)
	strs := strings.Split(str, "-")
	go func() {
		defer func() {
			close(ch)
		}()
		for _, str := range strs {
			select {
			case ch <- str:
			}
		}
	}()
	return nil
}
