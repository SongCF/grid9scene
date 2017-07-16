package main

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/SongCF/scene/gateway"
	. "github.com/SongCF/scene/model"
	. "github.com/SongCF/scene/util"
	"os"
	"github.com/SongCF/scene/rpc"
)

var server_started = make(chan struct{})

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

	//init rpc
	rpcAddr, err := Conf.Get(SCT_TCP, "rpc_server")
	CheckError(err)
	go rpc.InitServer(rpcAddr)

	go HttpServer()

	go TcpServer()

	go StartPProf()
	go StartStats()

	// register zk
	initZK()

	close(server_started)
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
