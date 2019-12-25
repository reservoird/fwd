package main

type ibooladapt interface {
	value() bool
}

type booladapt struct {
}

func (o *booladapt) value() bool {
	return true
}

// Fwd digester
type fwd struct {
	keepRunning ibooladapt
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	o.keepRunning = &booladapt{}
	return nil
}

// Digest reads from src channel and forwards to dst channel
func (o *fwd) Digest(src <-chan []byte, dst chan<- []byte) error {
	for o.keepRunning.value() == true {
		line := <-src
		dst <- line
	}
	return nil
}

// Digester for fwd
var Digester fwd
