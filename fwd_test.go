package main

import (
	"testing"
)

type flagtst struct {
	count int
	flag  bool
}

func newflagtst() *flagtst {
	f := new(flagtst)
	f.count = 0
	f.flag = true
	return f
}

func (o *flagtst) run() bool {
	val := o.flag
	if o.count == 0 {
		o.flag = false
	}
	return val
}

func TestConfig(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	if f.runner.run() == false {
		t.Errorf("expecting true but got false")
	}
}

func TestDigest(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	f.runner = newflagtst()
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
