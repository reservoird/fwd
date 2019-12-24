package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	if f.run() == false {
		t.Errorf("expecting true but got false")
	}
}

var one bool = true

func once() bool {
	val := one
	if one == true {
		one = false
	}
	return val
}

func TestDigest(t *testing.T) {
	f := fwd{}
	f.run = once
	src := make(chan []byte, 1)
	expected := []byte("hello")
	src <- expected
	dst := make(chan []byte, 1)
	err := f.Digest(src, dst)
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	actual := <-dst
	if string(actual) != string(expected) {
		t.Errorf("expecting %s but got %s", string(expected), string(actual))
	}

}
