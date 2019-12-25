package main

//go:generate mkdir -p mocks/iboolmocks
//go:generate mockgen -source $GOPATH/pkg/mod/github.com/reservoird/ibool@v1.0.0/ibool.go -package iboolmocks -destination mocks/iboolmocks/ibool_mock.go

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reservoird/fwd/mocks/iboolmocks"
)

func TestConfig(t *testing.T) {
	f := fwd{}
	err := f.Config("")
	if err != nil {
		t.Errorf("expecting nil but got error: %v", err)
	}
	if f.keepRunning.Val() == false {
		t.Errorf("expecting true but got false")
	}
}

func TestDigest(t *testing.T) {
	mockControl := gomock.NewController(t)
	defer mockControl.Finish()

	bmock := iboolmocks.NewMockIBool(mockControl)
	itr := 0
	itrTotal := 1
	bmock.EXPECT().Val().DoAndReturn(
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
	f.keepRunning = bmock
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
