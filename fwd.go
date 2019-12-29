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
	Tag       string
	Timestamp bool
}

// New is what reservoird to create and start fwd
func New(cfg string) (icd.Digester, error) {
	o := &fwd{
		Tag:       "fwd",
		Timestamp: false,
	}
	if cfg != "" {
		d, err := ioutil.ReadFile(cfg)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(d, o)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

// Name returns the digester name
func (o *fwd) Name() string {
	return o.Tag
}

// Digest reads from in queue and forwards to out queue
func (o *fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	o.run = true
	for o.run == true {
		d, err := iq.Get()
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			if d != nil {
				data, ok := d.([]byte)
				if ok == false {
					fmt.Printf("error invalid type\n")
				} else {
					line := string(data)
					if o.Timestamp == true {
						line = fmt.Sprintf("[%s %s] ", o.Name(), time.Now().Format(time.RFC3339)) + line
					}
					err = oq.Put([]byte(line))
					if err != nil {
						fmt.Printf("%v\n", err)
					}
				}
			}
		}

		select {
		case <-done:
			o.run = false
		default:
			time.Sleep(time.Millisecond)
		}
	}
	return nil
}
