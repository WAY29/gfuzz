package encoders

import (
	"encoding/base64"
)

func init() {
	AddEncoder("base64", &EncoderBase64{})
	AddEncoderInfo("base64", "Encodes the given string using base64.")
}

type EncoderBase64 struct {
}

func (p *EncoderBase64) Encode(s interface{}) interface{} {
	res := base64.StdEncoding.EncodeToString([]byte(s.(string)))
	return res
}
