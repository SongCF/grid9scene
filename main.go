package main

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/gateway"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	. "jhqc.com/songcf/scene/util"
	"os"
)

var server_started bool = false

func main() {
	defer RecoverPanic()

	log.SetLevel(log.DebugLevel)

	//
	wd, err := os.Getwd()
	CheckError(err)
	log.Println("work dir: ", wd)

	go HandleSignal()

	//load config
	InitConf()

	InitDB()
	InitCache()

	go HttpServer()

	go TcpServer()

	go StartPProf()
	go StartStats()

	// register zk
	initZK()

	server_started = true
	log.Info("server started!")

	select {}

	// 3.unregister zookeeper
	Unregister()
}

func initZK() {
	Register()
	GetServices("tcp")
	GetServices("http")
}
