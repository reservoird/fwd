package main

type runner interface {
	run() bool
}

type flag struct {
}

func (o *flag) run() bool {
	return true
}

// Fwd digester
type fwd struct {
	runner runner
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	o.runner = &flag{}
	return nil
}

// Digest reads from src channel and forwards to dst channel
func (o *fwd) Digest(src <-chan []byte, dst chan<- []byte) error {
	for o.runner.run() == true {
		line := <-src
		dst <- line
	}
	return nil
}

// Digester for fwd
var Digester fwd
