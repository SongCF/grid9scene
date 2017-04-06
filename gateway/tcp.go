package gateway

import (
	"net"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/util"
	"strconv"
	"time"
)


// will change by config file when start app
var (
	ReadBufSize = 2048
	WriteBufSize = 2048
	ReadDeadline time.Duration = 120 //second
)

func TcpServer() {
	tcpAddr := Conf.Get(SCT_TCP, "tcp_server")
	addr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	CheckError(err)

	listener, err := net.ListenTCP("tcp", addr)
	CheckError(err)
	log.Println("tcp server listening on: ", listener.Addr())

	loadTcpConfig()
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
		case <- GlobalDie:
			listener.Close()
			return
		default:
		}
	}
}

func loadTcpConfig() {
	//read write buf
	rBufSize, err := strconv.Atoi(Conf.Get(SCT_TCP, "read_buf"))
	CheckError(err)
	ReadBufSize = rBufSize
	wBufSize, err := strconv.Atoi(Conf.Get(SCT_TCP, "write_buf"))
	CheckError(err)
	WriteBufSize = wBufSize
	//deadline
	deadline, err := strconv.Atoi(Conf.Get(SCT_TCP, "deadline_time"))
	CheckError(err)
	ReadDeadline = time.Duration(deadline)
}
