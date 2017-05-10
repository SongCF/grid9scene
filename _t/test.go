package main

import (
	"fmt"
	"jhqc.com/songcf/scene/_test"
	"jhqc.com/songcf/scene/util"
	"net/http"
)

const (
	tcp_server   = ":9901"
	http_server  = "http://127.0.0.1:9911"
	min_interval = 0.3
	max_interval = 3.0
)

func main() {
	go util.HandleSignal()
	go func() {
		fmt.Println(http.ListenAndServe(":9999", nil))
	}()
	_test.TCPStressTest(http_server, tcp_server, min_interval, max_interval)
	select {}
}
