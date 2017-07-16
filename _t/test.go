package main

import (
	"fmt"
	"github.com/SongCF/scene/_test"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

const (
	tcp_server   = ":9901"
	http_server  = "http://127.0.0.1:9911"
	min_interval = 0.3
	max_interval = 3.0
)

func main() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		fmt.Println("Got signal:", s)
		os.Exit(1)
	}()
	go func() {
		fmt.Println(http.ListenAndServe(":9999", nil))
	}()
	_test.TCPStressTest(http_server, tcp_server, min_interval, max_interval)

	select {}
}
