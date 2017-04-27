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
	for {
		if server_started {
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	// start test client
	_test.StartClient(t)
}
