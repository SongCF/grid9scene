package main

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/gateway"
	. "jhqc.com/songcf/scene/model"
)



func main() {
	log.SetLevel(log.DebugLevel)

	go HandleSignal()

	InitDB()

	go HttpServer()

	go TcpServer()

	initZK()

	select {}
}



func initZK() {
	// TODO zookeeper
}