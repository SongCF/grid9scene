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

	writeBuf := make([]byte, 2048)
	for {
		select {
		case data, ok := <- s.ChanOut:
			if !ok {
				return
			}
			sendData(s.Conn, data, writeBuf)
		case <- GlobalDie:
			s.Conn.Close()
			//Don't close s.out, 如果grid server后退出，可能还会往里面写数据
			return
		}
	}
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