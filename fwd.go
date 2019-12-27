package main

import (
	"fmt"
	"sync"

	"github.com/reservoird/icd"
)

type fwd struct {
	run bool
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	return nil
}

// Digest reads from in queue and forwards to out queue
func (o *fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	fmt.Printf("fwd.digest: into\n")
	o.run = true
	for o.run == true {
		fmt.Printf("fwd.digest: before pop\n")
		d, err := iq.Pop()
		fmt.Printf("fwd.digest: after pop\n")
		if err != nil {
			return err
		}
		fmt.Printf("fwd.digest: before push\n")
		err = oq.Push(d)
		fmt.Printf("fwd.digest: after push\n")
		if err != nil {
			return err
		}

		select {
		case <-done:
			o.run = false
		default:
		}
	}
	fmt.Printf("fwd.digest: outof\n")
	return nil
}

// Digester for fwd
var Digester fwd
