package main

import (
	"testing"
)

type boolbridgetest struct {
	count int
	val   bool
}

func newboolbridgetest() *boolbridgetest {
	b := new(boolbridgetest)
	b.count = 0
	b.val = true
	return b
}

func (o *boolbridgetest) value() bool {
	val := o.val
	if o.count == 0 {
		o.val = false
	}
	o.count = o.count + 1
	return val
}

func TestConfig(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	if f.keepRunning.value() == false {
		t.Errorf("expecting true but got false")
	}
}

func TestDigest(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	f.keepRunning = newboolbridgetest()
	src := make(chan []byte, 1)
	expected := []byte("hello")
	src <- expected
	dst := make(chan []byte, 1)
	err = f.Digest(src, dst)
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	actual := <-dst
	if string(actual) != string(expected) {
		t.Errorf("expecting %s but got %s", string(expected), string(actual))
	}

}
