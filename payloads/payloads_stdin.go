package payloads

import (
	"bufio"
	"os"
)

func init() {
	AddPayload("stdin", &PayloadStdin{})
	AddPayloadInfo("stdin", "Returns each line from stdin.")
}

type PayloadStdin struct {
	channel chan interface{}
}

//
func (p *PayloadStdin) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadStdin) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	br := bufio.NewReader(os.Stdin)
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			bytes, _, err := br.ReadLine()
			line := string(bytes[:])
			if err != nil {
				return
			}
			select {
			case ch <- line:
				continue
			}
		}
	}()
	return nil
}
