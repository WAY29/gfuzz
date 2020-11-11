package encoders

func init() {
	AddEncoder("example", &EncoderExample{})
	AddEncoderInfo("example", "Check if you want to write a custom encoder (This is Introduction of a encoder)")
}

type EncoderExample struct {
}

func (p *EncoderExample) Encode(s interface{}) interface{} {
	return s.(string)
}
