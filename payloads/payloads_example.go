package payloads

func init() {
	AddPayload("example", &PayloadExample{})
	AddPayloadInfo("example", "Check if you want to write a custom payload (This is Introduction of a payload)")
}

type PayloadExample struct {
	channel chan interface{}
}

//
func (p *PayloadExample) Channel() chan interface{} {
	return p.channel
}

func (p *PayloadExample) New(s ...interface{}) error {
	p.channel = make(chan interface{})
	ch := p.channel
	str := s[0].(string)
	strs := []string{str}
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
