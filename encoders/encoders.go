package encoders

var EncodersArray map[string]Encoder = make(map[string]Encoder)
var EncodersInfoArray map[string]string = make(map[string]string)

// * Pyaloads encoder interface, use go doc encoders.EncoderBase for more information
type Encoder interface {
	Encode(s interface{}) interface{}
}

/*
* Encoder base struct, must implement Encode method

* To write a custom encoder, please check encoders_example.go
 */
type EncoderBase struct {
}

func (p *EncoderBase) Encode(s interface{}) interface{} {
	return s
}

func AddEncoder(name string, e Encoder) {
	EncodersArray[name] = e
}

func AddEncoderInfo(name string, info string) {
	EncodersInfoArray[name] = info
}

func GetEncoder(name string) Encoder {
	e, ok := EncodersArray[name]
	if !ok {
		return &EncoderBase{}
	} else {
		return e
	}
}

func GetEncoderInfo(name string) string {
	s, ok := EncodersInfoArray[name]
	if !ok {
		return ""
	} else {
		return s
	}
}
