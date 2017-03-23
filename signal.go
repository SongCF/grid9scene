package main

import (
	"sync"
	"os"
	"os/signal"
	"syscall"
	log "github.com/Sirupsen/logrus"
)


var (
	globalWG sync.WaitGroup
	// server close signal
	globalDie = make(chan struct{})
)

func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	for {
		msg := <- ch
		switch msg {
		case syscall.SIGTERM:
			close(globalDie)
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")

			// waiting global server end.
			// 1.tcp listener
			// 2.session reader/sender/agent/conn
			// 3.grid server
			// 4. TODO unregister zookeeper
			globalWG.Wait()
			log.Info("shutdown.")
			os.Exit(0)
		}
	}
}