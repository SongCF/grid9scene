package main

import (
	"jhqc.com/songcf/scene/_test"
	"testing"
	"time"
)

func TestScene(t *testing.T) {
	// start server
	go main()
	// waiting started
	select {
	case <-server_started:
	case <-time.After(time.Second * 60):
		t.Fatal("start server timeout!")
	}
	// start test client
	_test.StartClient(t)
}
