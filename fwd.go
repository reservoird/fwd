package main

import (
	"github.com/reservoird/ibool"
)

type boolbridge struct {
}

func (o *boolbridge) Val() bool {
	return true
}

// Fwd digester
type fwd struct {
	keepRunning ibool.IBool
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	o.keepRunning = &boolbridge{}
	return nil
}

// Digest reads from src channel and forwards to dst channel
func (o *fwd) Digest(src <-chan []byte, dst chan<- []byte) error {
	for o.keepRunning.Val() == true {
		line := <-src
		dst <- line
	}
	return nil
}

// Digester for fwd
var Digester fwd
