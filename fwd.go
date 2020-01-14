package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
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
}

// Fwd contains what is needed to run digester
type Fwd struct {
	cfg FwdCfg
	run bool
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
		cfg: c,
		run: false,
	}
	return o, nil
}

// Name returns the digester name
func (o *Fwd) Name() string {
	return o.cfg.Name
}

// Running returns whether or not digest is running
func (o *Fwd) Running() bool {
	return o.run
}

// Digest reads from rcv queue and forwards to snd queue
func (o *Fwd) Digest(rcv icd.Queue, snd icd.Queue, mc *icd.MonitorControl) {
	defer mc.WaitGroup.Done()

	stats := FwdStats{}

	o.run = true
	stats.Name = o.cfg.Name
	stats.Running = o.run
	for o.run == true {
		if rcv.Closed() == false {
			d, err := rcv.Get()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				if d != nil {
					data, ok := d.([]byte)
					if ok == false {
						fmt.Printf("error invalid type\n")
					} else {
						stats.MessagesReceived = stats.MessagesReceived + 1
						if snd.Closed() == false {
							line := string(data)
							if o.cfg.Timestamp == true {
								line = fmt.Sprintf("[%s %s] ", o.Name(), time.Now().Format(time.RFC3339)) + line
							}
							err = snd.Put([]byte(line))
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
		case <-mc.ClearChan:
			stats = FwdStats{
				Name:    o.cfg.Name,
				Running: o.run,
			}
		default:
		}

		// send to monitor
		select {
		case mc.StatsChan <- stats:
		default:
		}

		// listens for shutdown
		select {
		case <-mc.DoneChan:
			o.run = false
			stats.Running = o.run
		default:
		}

		runtime.Gosched()
	}

	// send final stats blocking
	mc.FinalStatsChan <- stats
}
