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
	Name             string
	MessagesReceived uint64
	MessagesSent     uint64
	Running          bool
	Monitoring       bool
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
		run:       false,
		statsChan: make(chan FwdStats, 1),
		clearChan: make(chan struct{}, 1),
	}
	return o, nil
}

// Name returns the digester name
func (o *Fwd) Name() string {
	return o.cfg.Name
}

// Monitor provides stats and clear
func (o *Fwd) Monitor(statsChan chan<- string, clearChan <-chan struct{}, doneChan <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done() // required

	stats := FwdStats{}
	monrun := true
	for monrun == true {
		// clear
		select {
		case <-clearChan:
			select {
			case o.clearChan <- struct{}{}:
			default:
			}
		default:
		}

		// done
		select {
		case <-doneChan:
			monrun = false
			stats.Monitoring = monrun
		default:
		}

		// get stats from digest
		select {
		case stats = <-o.statsChan:
			stats.Monitoring = monrun
		default:
		}

		// marshal
		data, err := json.Marshal(stats)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			// send stats to reservoird
			select {
			case statsChan <- string(data):
			default:
			}
		}

		if monrun == true {
			time.Sleep(time.Millisecond)
		}
	}
}

// Running returns whether or not digest is running
func (o *Fwd) Running() bool {
	return o.run
}

// Digest reads from in queue and forwards to out queue
func (o *Fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	stats := FwdStats{}

	o.run = true
	stats.Name = o.cfg.Name
	stats.Running = o.run
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

		// clear
		select {
		case <-o.clearChan:
			stats = FwdStats{}
		default:
		}

		// listens for shutdown
		select {
		case <-done:
			o.run = false
			stats.Name = o.cfg.Name
			stats.Running = o.run
		default:
		}

		// send to monitor
		select {
		case o.statsChan <- stats:
		default:
		}

		if o.run == true {
			time.Sleep(time.Millisecond)
		}
	}
	return nil
}
