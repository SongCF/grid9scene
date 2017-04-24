package gateway

import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	"io"
	. "jhqc.com/songcf/scene/controller"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	. "jhqc.com/songcf/scene/util"
	"net"
	"time"
)

func handleClient(conn net.Conn) {
	defer log.Debug("---session reader end.")
	defer RecoverPanic()

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
	s.Die = make(chan struct{})

	defer CloseSession(&s)

	in := make(chan []byte)
	defer close(in)

	go agent(&s, in)

	GlobalWG.Add(1)
	go sender(&s)

	//reader
	readBuf := make([]byte, ReadBufSize)
	for {
		conn.SetReadDeadline(time.Now().Add(ReadDeadline))

		n, err := io.ReadAtLeast(conn, readBuf[:4], 4)
		if err != nil {
			log.Info("read header error:", err)
			return
		}
		size := binary.BigEndian.Uint32(readBuf[:4])

		// read data
		n, err = io.ReadAtLeast(conn, readBuf[:size], int(size))
		if err != nil {
			log.Errorf("read payload failed, ip:%v reason:%v size:%v\n", s.IP, err, n)
			return
		}

		select {
		case in <- readBuf[:size]:
		case <-s.Die:
			return
		case <-GlobalDie:
			return
		}
	}
}
