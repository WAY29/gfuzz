package payloads

import (
	"bufio"
	"gfuzz/utils"
	"os"
)

func init() {
	AddPayload("file", &PayloadFile{})
	AddPayloadInfo("file", "Returns each word from a file.")
}

type PayloadFile struct {
	channel chan interface{}
}

//
func (p *PayloadFile) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadFile) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	f, err := os.Open(s[0].(string))
	if err != nil {
		utils.PrintError("No such file")
		return err
	}
	br := bufio.NewReader(f)
	go func() {
		defer func() {
			f.Close()
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
