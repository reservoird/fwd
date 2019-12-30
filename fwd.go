package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/reservoird/icd"
)

// FwdCfg contains fwd config
type FwdCfg struct {
	Name      string
	Timestamp bool
}

// FwdStats contains fwd stats
type FwdStats struct {
	MessagesReceived uint64
	MessagesSent     uint64
	Running          bool
}

// Fwd contains what is needed to run digester
type Fwd struct {
	cfg       FwdCfg
	run       bool
	statsChan chan FwdStats
	clearChan chan struct{}
}

// New is what reservoird to create and start fwd
func New(cfg string) (icd.Digester, error) {
	c := FwdCfg{
		Name:      "com.reservoird.digest.fwd",
		Timestamp: false,
	}
	if cfg != "" {
		d, err := ioutil.ReadFile(cfg)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(d, &c)
		if err != nil {
			return nil, err
		}
	}
	o := &Fwd{
		cfg:       c,
		run:       true,
		statsChan: make(chan FwdStats),
		clearChan: make(chan struct{}),
	}
	return o, nil
}

// Name returns the digester name
func (o *Fwd) Name() string {
	return o.cfg.Name
}

// Stats returns stats NOTE: thread safe
func (o *Fwd) Stats() (string, error) {
	select {
	case stats := <-o.statsChan:
		data, err := json.Marshal(stats)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
	}
	return "", nil
}

// ClearStats clears stats NOTE: thread safe
func (o *Fwd) ClearStats() {
	select {
	case o.clearChan <- struct{}{}:
	default:
	}
}

// Digest reads from in queue and forwards to out queue
func (o *Fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	stats := FwdStats{}

	o.run = true
	for o.run == true {
		if iq.Closed() == false {
			d, err := iq.Get()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				if d != nil {
					data, ok := d.([]byte)
					if ok == false {
						fmt.Printf("error invalid type\n")
					} else {
						stats.MessagesReceived = stats.MessagesReceived + 1
						if oq.Closed() == false {
							line := string(data)
							if o.cfg.Timestamp == true {
								line = fmt.Sprintf("[%s %s] ", o.Name(), time.Now().Format(time.RFC3339)) + line
							}
							err = oq.Put([]byte(line))
							if err != nil {
								fmt.Printf("%v\n", err)
							} else {
								stats.MessagesSent = stats.MessagesSent + 1
							}
						}
					}
				}
			}
		}

		select {
		case <-o.clearChan:
			stats = FwdStats{}
		default:
		}

		select {
		case <-done:
			o.run = false
		default:
		}

		stats.Running = o.run
		select {
		case o.statsChan <- stats:
		default:
		}

		time.Sleep(time.Millisecond)
	}
	return nil
}
