package gateway

import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	"io"
	. "jhqc.com/songcf/scene/controller"
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
	const lenHead = 4
	head := make([]byte, lenHead)
	for {
		conn.SetReadDeadline(time.Now().Add(ReadDeadline))

		n, err := io.ReadFull(conn, head)
		if err != nil {
			log.Info("read header error:", err)
			return
		}
		size := binary.BigEndian.Uint32(head)

		// read data
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)
		if err != nil {
			log.Errorf("read payload failed, ip:%v reason:%v size:%v\n", s.IP, err, n)
			return
		}

		select {
		case in <- data:
		case <-s.Die:
			return
		case <-GlobalDie:
			return
		}
	}
}
