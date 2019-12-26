package main

import "sync"

type fwd struct {
	run bool
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	return nil
}

// Digest reads from src channel and forwards to dst channel
func (o *fwd) Digest(src <-chan []byte, dst chan<- []byte, done <-chan struct{}, wg *sync.WaitGroup) error {
	for o.run == true {
		line := <-src
		dst <- line
		select {
		case <-done:
			o.run = false
		default:
		}
	}
	return nil
}

// Digester for fwd
var Digester fwd
