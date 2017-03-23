package main

import (
	"net"
	log "github.com/Sirupsen/logrus"
)

func tcpServer() {
	addr, err := net.ResolveTCPAddr("tcp4", ":9901")
	checkError(err)

	listener, err := net.ListenTCP("tcp", addr)
	checkError(err)
	log.Info("listening on:", listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warning("accept failed:", err)
			continue
		}

		// set socket read buffer
		conn.SetReadBuffer(256)   //TODO config
		// set socket write buffer
		conn.SetWriteBuffer(2048) //TODO config
		// start a goroutine for every incoming connection for reading
		go handleClient(conn)

		// check server close signal
		select {
		case <- globalDie:
			listener.Close()
			return
		default:
		}
	}
}
