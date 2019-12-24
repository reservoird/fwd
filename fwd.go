package main

// Fwd digester
type fwd struct {
	run func() bool
}

func forever() bool {
	return true
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	o.run = forever
	return nil
}

// Digest reads from src channel and forwards to dst channel
func (o *fwd) Digest(src <-chan []byte, dst chan<- []byte) error {
	for o.run() {
		line := <-src
		dst <- line
	}
	return nil
}

// Digester for fwd
var Digester fwd
