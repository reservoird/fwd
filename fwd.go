package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/reservoird/icd"
)

type fwdCfg struct {
	Name      string
	Timestamp bool
}

type fwdStats struct {
	MsgsRead uint64
	MsgsDrop uint64
	MsgsSent uint64
	Active   bool
}

type fwd struct {
	cfg       fwdCfg
	stats     fwdStats
	statsChan chan<- string
}

// New is what reservoird to create and start fwd
func New(cfg string, statsChan chan<- string) (icd.Digester, error) {
	c := fwdCfg{
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
	o := &fwd{
		cfg:       c,
		stats:     fwdStats{},
		statsChan: statsChan,
	}
	return o, nil
}

// Name returns the digester name
func (o *fwd) Name() string {
	return o.cfg.Name
}

// Digest reads from in queue and forwards to out queue
func (o *fwd) Digest(iq icd.Queue, oq icd.Queue, done <-chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()

	o.stats.Active = true
	for o.stats.Active == true {
		if iq.Closed() == false {
			d, err := iq.Get()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				if d != nil {
					data, ok := d.([]byte)
					if ok == false {
						fmt.Printf("error invalid type\n")
						o.stats.MsgsDrop = o.stats.MsgsDrop + 1
					} else {
						o.stats.MsgsRead = o.stats.MsgsRead + 1
						if oq.Closed() == false {
							line := string(data)
							if o.cfg.Timestamp == true {
								line = fmt.Sprintf("[%s %s] ", o.Name(), time.Now().Format(time.RFC3339)) + line
							}
							err = oq.Put([]byte(line))
							if err != nil {
								fmt.Printf("%v\n", err)
								o.stats.MsgsDrop = o.stats.MsgsDrop + 1
							} else {
								o.stats.MsgsSent = o.stats.MsgsSent + 1
							}
						}
					}
				} else {
					o.stats.MsgsDrop = o.stats.MsgsDrop + 1
				}
			}
		}

		select {
		case <-done:
			o.stats.Active = false
		default:
		}

		stats, err := json.Marshal(o.stats)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		select {
		case o.statsChan <- string(stats):
		default:
		}

		time.Sleep(time.Millisecond)
	}
	return nil
}
