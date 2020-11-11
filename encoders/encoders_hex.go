package encoders

import (
	"encoding/hex"
)

func init() {
	AddEncoder("hex", &EncoderHex{})
	AddEncoderInfo("hex", "Every byte of data is converted into the corresponding 2-digit hex representation.")
}

type EncoderHex struct {
}

func (p *EncoderHex) Encode(s interface{}) interface{} {
	return hex.EncodeToString([]byte(s.(string)))
}
