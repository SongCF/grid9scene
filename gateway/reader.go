package gateway

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"io"
	"net"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	"encoding/binary"
)

func handleClient(conn net.Conn) {
	defer func() {
		log.Debug("---session reader end.")
	}()

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("init session error:", err)
		return
	}
	log.Infof("new connection from: %v port:%v", host, port)

	// init session
	var s Session
	s.IP = net.ParseIP(host)
	s.Conn = conn
	s.PacketCount = 0
	s.ConnectTime = time.Now()
	s.ChanOut = make(chan []byte)


	in := make(chan []byte)
	defer close(in)

	go agent(&s, in)

	GlobalWG.Add(1)
	go sender(&s)

	//reader
	readBuf := make([]byte, 256)   // TODO config size
	for {
		readDeadline := 120 * time.Second  //TODO config
		conn.SetReadDeadline(time.Now().Add(readDeadline))

		n, err := io.ReadAtLeast(conn, readBuf[:4], 4)
		if err != nil {
			log.Error("read header error:", err)
			return
		}
		size := binary.BigEndian.Uint32(readBuf[:4])

		// read data
		n, err = io.ReadAtLeast(conn, readBuf[:size], int(size))
		if err != nil {
			log.Errorf("read payload failed, ip:%v reason:%v size:%v\n", s.IP, err, n)
			return
		}
		log.Debugf("---read head size=%v, n=%v", size, n)

		select {
		case in <- readBuf[:size]:
		case <- s.Die:
			return
		case <- GlobalDie:
			return
		}
	}
}

