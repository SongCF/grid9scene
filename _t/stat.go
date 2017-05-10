package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	stat_server = "http://127.0.0.1:9912"
)

func main() {
	displayStat()
}

func displayStat() {
	for {
		fmt.Println("displayStat...")
		uri := fmt.Sprintf("%s/stat/avg_msg", stat_server)
		resp, err := http.Get(uri)
		if err != nil {
			fmt.Printf("displayStat req error:%v\n", err)
			time.Sleep(time.Second * 10)
			continue
		}
		rBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("displayStat read body error:%v\n", err)
			time.Sleep(time.Second * 10)
			continue
		}
		fmt.Printf("%s\n", string(rBody))
		resp.Body.Close()
		time.Sleep(time.Second * 10)
	}
}
