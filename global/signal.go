package global

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	GlobalWG sync.WaitGroup
	// server close signal
	GlobalDie = make(chan struct{})
)

func HandleSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM:
			close(GlobalDie)
			log.Info("sigterm received")
			log.Info("waiting for agents close, please wait...")

			// waiting global server end.
			// 1.tcp listener
			// 2.session reader/sender/agent/conn
			// 3.app/space/grid server
			// 4. TODO unregister zookeeper
			GlobalWG.Wait()
			log.Info("shutdown.")
			os.Exit(0)
		}
	}
}
