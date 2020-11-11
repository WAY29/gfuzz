package encoders

import "testing"

func TestEncoders(t *testing.T) {
	EncodersStringArray := []string{"base64", "md5", "hex"}
	InputArray := []string{"abc", "admin", "qwe"}
	ExpectedArray := []string{"YWJj", "21232f297a57a5a743894a0e4a801fc3", "717765"}
	for i, s := range EncodersStringArray {
		e := GetEncoder(s)
		actual := e.Encode(InputArray[i])
		if actual != ExpectedArray[i] {
			t.Fatalf("Encoder %s expected = %s, actual: %s", s, ExpectedArray[i], actual)
		}
	}
}
