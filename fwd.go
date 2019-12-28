package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/reservoird/icd"
)

type fwd struct {
	run       bool
	Name      string
	Timestamp bool
}

// NewDigester is what reservoird to create and start fwd
func NewDigester() (icd.Digester, error) {
	return new(fwd), nil
}

// Config configures digester
func (o *fwd) Config(cfg string) error {
	d, err := ioutil.ReadFile(cfg)
	if err != nil {
		return err
	}
	f := fwd{}
	err = json.Unmarshal(d, &f)
	if err != nil {
		return err
	}
	o.Name = f.Name
	o.Timestamp = f.Timestamp
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
		data, ok := d.([]byte)
		if ok == false {
			return fmt.Errorf("error invalid type")
		}
		line := string(data)
		if o.Timestamp == true {
			line = fmt.Sprintf("%s %s: ", o.Name, time.Now().Format(time.RFC3339)) + line
		}
		err = oq.Push([]byte(line))
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
