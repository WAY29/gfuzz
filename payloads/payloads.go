package payloads

var PayloadsArray map[string]Payload = make(map[string]Payload)
var PayloadsInfoArray map[string]string = make(map[string]string)

// * Pyaloads interface, use go doc payloads.PyaloadsBase for more information
type Payload interface {
	Channel() chan interface{}
	New(s ...interface{}) error
}

/*
* Payload base struct, must implement Channel and New method

* To write a custom payload, please check payloads_example.go
 */
type PayloadBase struct {
	channel chan interface{}
}

func AddPayload(name string, e Payload) {
	PayloadsArray[name] = e
}

func AddPayloadInfo(name string, info string) {
	PayloadsInfoArray[name] = info
}

func GetPayload(name string) Payload {
	e, ok := PayloadsArray[name]
	if !ok {
		return nil
	} else {
		return e
	}
}

func GetPayloadInfo(name string) string {
	s, ok := PayloadsInfoArray[name]
	if !ok {
		return ""
	} else {
		return s
	}
}
