package gateway

import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	. "jhqc.com/songcf/scene/util"
	"net"
)

func sender(s *Session) {
	defer log.Debug("---session sender end.")
	defer GlobalWG.Done()
	defer RecoverPanic()
	defer s.Close()
	defer s.Conn.Close()

	writeBuf := make([]byte, WriteBufSize)
	for {
		select {
		case data, ok := <-s.ChanOut:
			if !ok {
				return
			}
			sendData(s.Conn, data, writeBuf)
		case <-s.Die:
			return
		case <-GlobalDie:
			return
		}
	}
}

func sendData(conn net.Conn, data []byte, cache []byte) bool {
	size := len(data)
	binary.BigEndian.PutUint32(cache, uint32(size))
	copy(cache[4:], data)

	n, err := conn.Write(cache[:size+4])
	if err != nil {
		log.Errorf("... send data error, bytes:%v, reason:%v", n, err)
		return false
	}
	return true
}
