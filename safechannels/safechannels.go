package safechannels

type T bool

type SafeChannels struct {
	ch chan T
}

func New() *SafeChannels {
	return &SafeChannels{ch: make(chan T, 1)}
}

func (safech *SafeChannels) Channel() chan T {
	return safech.ch
}

func (safech *SafeChannels) IsClosed() bool {
	select {
	case <-safech.ch:
		return true
	default:
	}
	return false
}

func (safech *SafeChannels) SafeClose() (justClosed bool) {
	defer func() {
		if recover() != nil {
			// The return result can be altered
			// in a defer function call.
			justClosed = false
		}
	}()

	// assume ch != nil here.
	close(safech.ch) // panic if ch is closed
	return true      // <=> justClosed = true; return
}

func (safech *SafeChannels) SafeSend(value T) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()

	safech.ch <- value // panic if ch is closed
	return false       // <=> closed = false; return
}
