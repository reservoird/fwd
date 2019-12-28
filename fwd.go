package main

import (
	"sync"

	"github.com/reservoird/icd"
)

type fwd struct {
	run bool
}

// NewDigester is what reservoird to create and start fwd
func NewDigester() (icd.Digester, error) {
	return new(fwd), nil
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	return nil
}

// Digest reads from in queue and forwards to out queue
func (o *fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	o.run = true
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
