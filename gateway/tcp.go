package gateway

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/util"
	"net"
)

func TcpServer() {
	tcpAddr, err := Conf.Get(SCT_TCP, "tcp_server")
	CheckError(err)
	addr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	CheckError(err)

	listener, err := net.ListenTCP("tcp", addr)
	CheckError(err)
	log.Println("tcp server listening on: ", listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warning("accept failed:", err)
			continue
		}

		// set socket read buffer
		conn.SetReadBuffer(ReadBufSize)
		// set socket write buffer
		conn.SetWriteBuffer(WriteBufSize)
		// start a goroutine for every incoming connection for reading
		go handleClient(conn)

		// check server close signal
		select {
		case <-GlobalDie:
			listener.Close()
			return
		default:
		}
	}
}
