package main

import (
	"sync"

	"github.com/reservoird/reservoird/run"
)

type fwd struct {
	run bool
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	return nil
}

// Digest reads from in queue and forwards to out queue
func (o *fwd) Digest(iq run.Queue, oq run.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	for o.run == true {
		d, err := iq.Pop()
		if err != nil {
			return err
		}
		err = oq.Push(d)
		if err != nil {
			return err
		}

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
