package main

import (
	"time"
	"net"
	log "github.com/Sirupsen/logrus"
	"io"
	"encoding/binary"
	. "jhqc.com/songcf/scene/types"
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

	globalWG.Add(1)
	go sender(&s)

	//reader
	readBuf := make([]byte, 256)   // TODO config size
	for {
		readDeadline := 120 * time.Second  //TODO config
		conn.SetReadDeadline(time.Now().Add(readDeadline))

		n, err := io.ReadAtLeast(conn, readBuf, 4)
		if err != nil {
			log.Error("read header error:", err)
			return
		}
		size := binary.BigEndian.Uint32(readBuf[:4])

		// read data
		n, err = io.ReadAtLeast(conn, readBuf, int(size))
		if err != nil {
			log.Errorf("read payload failed, ip:%v reason:%v size:%v\n", s.IP, err, n)
			return
		}

		select {
		case in <- readBuf[:size]:
		case <- globalDie:
			return
		}
	}
}

func agent(s *Session, in chan []byte) {
	defer func() {
		log.Debug("---session agent end.")
	}()

	minTimer := time.After(time.Minute)
	for {
		select {
		case msg, ok := <- in:
			if !ok {
				return
			}
			s.PacketCount++
			log.Infof("req msg:%v, packCount:%v", msg, s.PacketCount)
			msgHandler(s.AppId, s.Uid, msg)
		case <- minTimer:
			timeWork()
			minTimer = time.After(time.Minute)
		case <- globalDie:
			return
		}
	}
}

func sender(s *Session) {
	defer func() {
		log.Debug("---session sender end.")
	}()
	defer globalWG.Done()

	writeBuf := make([]byte, 2048)
	for {
		select {
		case data, ok := <- s.ChanOut:
			if !ok {
				return
			}
			sendData(s.Conn, data, writeBuf)
		case <- globalDie:
			s.Conn.Close()
			//Don't close s.out, 如果grid server后退出，可能还会往里面写数据
			return
		}
	}
}





func timeWork() {
	//TODO something
	log.Info("on minute timer")
}

func sendData(conn net.Conn, data []byte, cache []byte) bool {
	log.Info("... send data ...")
	size := len(data)
	binary.BigEndian.PutUint32(cache, uint32(size))
	copy(cache[4:], data)  //TODO 4, config

	n, err := conn.Write(cache[:size+4])
	if err != nil {
		log.Errorf("... send data error, bytes:%v, reason:%v", n, err)
		return false
	}
	log.Info("... send data ok!")
	return true
}


func msgHandler(appId string, uid int32, m []byte) {
	cmd := 1
	switch cmd {
	case 1:

	case 2:
		//s := AppList[appId].sessionList[uid]
		//msg2Grid(s.appId, s.spaceId, s.gridId, m)
	default:
	}
}