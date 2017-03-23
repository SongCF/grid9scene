package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)


func main() {
	log.SetLevel(log.DebugLevel)

	go signalHandler()

	initDB()

	go httpServer()

	go tcpServer()

	initZK()

	select {}
}


func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func initZK() {
	// TODO zookeeper
}