package main

import (
	"time"
	"net"
	log "github.com/Sirupsen/logrus"
	"io"
	"encoding/binary"
)

type Session struct {
	ip           net.IP
	conn         net.Conn

	appId        string
	uid          int32

	in           chan []byte
	out          chan []byte
	ntf          chan []byte
	die          chan struct{} //会话关闭信号

	readBuf      []byte
	writeBuf     []byte

	packetCount int32         //对包进行计数
	connectTime time.Time
}

func (s* Session) init(conn net.Conn) bool {
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("init session error:", err)
		return false
	}
	log.Infof("new connection from: %v port:%v", host, port)

	s.ip = net.ParseIP(host)
	s.conn = conn

	s.die = make(chan struct{})
	s.in = make(chan []byte)
	s.out = make(chan []byte)
	s.ntf = make(chan []byte)

	s.readBuf = make([]byte, 256)  //TODO config
	s.writeBuf = make([]byte, 2048) //TODO config

	s.packetCount = 0
	s.connectTime = time.Now()

	return true
}



func handleClient(conn net.Conn) {
	// init session
	var s Session
	ok := s.init(conn)
	if !ok {
		log.Error("init seesion error.")
		return
	}
	defer close(s.in) //客户端conn关闭时，read header失败，session退出，在这里关闭 s.in来触发agent结束

	go agent(&s)
	go sender(&s)

	//reader
	header := make([]byte, 4)
	for {
		// solve dead link problem:
		// physical disconnection without any communication between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		readDeadline := 120 * time.Second  //TODO config
		conn.SetReadDeadline(time.Now().Add(readDeadline))

		n, err := io.ReadFull(conn, header)
		if err != nil {
			log.Error("read header error:", err)
			return
		}
		size := binary.BigEndian.Uint32(n)

		// read data
		n, err = io.ReadAtLeast(conn, s.readBuf, size)
		if err != nil {
			log.Errorf("read payload failed, ip:%v reason:%v size:%v\n", s.ip, err, n)
			return
		}

		select {
		case s.in <- s.readBuf:
		//case <- global_die:
		}
	}
}

func agent(s *Session) {

}

func sender(s *Session) {

}