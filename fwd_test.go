package main

import (
	"testing"
)

type booladapttst struct {
	count int
	val   bool
}

func newbooladapttst() *booladapttst {
	f := new(booladapttst)
	f.count = 0
	f.val = true
	return f
}

func (o *booladapttst) value() bool {
	val := o.val
	if o.count == 0 {
		o.val = false
	}
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
	f.keepRunning = newbooladapttst()
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
