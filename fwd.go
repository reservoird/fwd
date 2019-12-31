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
	stats     FwdStats
	statsLock sync.Mutex
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
		stats:     FwdStats{},
		statsLock: sync.Mutex{},
	}
	return o, nil
}

// Name returns the digester name
func (o *Fwd) Name() string {
	return o.cfg.Name
}

// Stats returns stats NOTE: thread safe
func (o *Fwd) Stats() (string, error) {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	data, err := json.Marshal(o.stats)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ClearStats clears stats NOTE: thread safe
func (o *Fwd) ClearStats() {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	o.stats = FwdStats{}
}

func (o *Fwd) incMessagesReceived() {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	o.stats.MessagesReceived = o.stats.MessagesReceived + 1
}

func (o *Fwd) incMessagesSent() {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	o.stats.MessagesSent = o.stats.MessagesSent + 1
}

func (o *Fwd) setRunning(run bool) {
	o.statsLock.Lock()
	defer o.statsLock.Unlock()
	o.stats.Running = run
}

// Digest reads from in queue and forwards to out queue
func (o *Fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	o.run = true
	o.setRunning(o.run)
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
						o.incMessagesReceived()
						if oq.Closed() == false {
							line := string(data)
							if o.cfg.Timestamp == true {
								line = fmt.Sprintf("[%s %s] ", o.Name(), time.Now().Format(time.RFC3339)) + line
							}
							err = oq.Put([]byte(line))
							if err != nil {
								fmt.Printf("%v\n", err)
							} else {
								o.incMessagesSent()
							}
						}
					}
				}
			}
		}

		select {
		case <-done:
			o.run = false
			o.setRunning(o.run)
		default:
		}

		time.Sleep(time.Millisecond)
	}
	return nil
}
