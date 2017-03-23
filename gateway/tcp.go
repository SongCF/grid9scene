package gateway

import (
	"net"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	"jhqc.com/songcf/scene/util"
)

func TcpServer() {
	addr, err := net.ResolveTCPAddr("tcp4", ":9901")
	util.CheckError(err)

	listener, err := net.ListenTCP("tcp", addr)
	util.CheckError(err)
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
		case <- GlobalDie:
			listener.Close()
			return
		default:
		}
	}
}
