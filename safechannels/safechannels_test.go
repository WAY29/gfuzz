package safechannels

import (
	"testing"
)

func TestSafeChannels(t *testing.T) {
	safech := New()
	safech.SafeSend(true)
	t.Logf("%#v", safech.IsClosed())
	t.Logf("%#v", <-safech.ch)
}
