package encoders

import (
	"crypto/md5"
	"encoding/hex"
)

func init() {
	AddEncoder("md5", &EncoderMd5{})
	AddEncoderInfo("md5", "Applies a md5 hash to the given string.")
}

type EncoderMd5 struct {
}

func (p *EncoderMd5) Encode(s interface{}) interface{} {
	e := md5.New()
	_, err := e.Write([]byte(s.(string)))
	if err == nil {
		return hex.EncodeToString(e.Sum(nil))
	}
	return ""
}
