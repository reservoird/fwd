package main

//go:generate mockgen -source fwd.go -package main -destination fwd_mocks.go

import (
	"testing"

	"github.com/golang/mock/gomock"
)

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
	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	iboolmock := NewMockibool(mockControl)
	itr := 0
	itrTotal := 1
	iboolmock.EXPECT().value().DoAndReturn(
		func() bool {
			if itr == itrTotal {
				return false
			}
			itr = itr + 1
			return true
		},
	).Times(2)

	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	f.keepRunning = iboolmock
	src := make(chan []byte, 2)
	expected := []byte("hello")
	src <- expected
	dst := make(chan []byte, 2)
	err = f.Digest(src, dst)
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	actual := <-dst
	if string(actual) != string(expected) {
		t.Errorf("expecting %s but got %s", string(expected), string(actual))
	}

}
