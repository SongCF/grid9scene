package gateway

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"encoding/binary"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
)

func sender(s *Session) {
	defer func() {
		log.Debug("---session sender end.")
	}()
	defer GlobalWG.Done()

	writeBuf := make([]byte, WriteBufSize)
	for {
		select {
		case data, ok := <- s.ChanOut:
			if !ok {
				return
			}
			sendData(s.Conn, data, writeBuf)
		case <- s.Die:
			s.Conn.Close()
			return
		case <- GlobalDie:
			s.Conn.Close()
			return
		}
	}
}

func sendData(conn net.Conn, data []byte, cache []byte) bool {
	log.Info("... send data ...")
	size := len(data)
	binary.BigEndian.PutUint32(cache, uint32(size))
	copy(cache[4:], data)

	n, err := conn.Write(cache[:size+4])
	if err != nil {
		log.Errorf("... send data error, bytes:%v, reason:%v", n, err)
		return false
	}
	log.Info("... send data ok!")
	return true
}