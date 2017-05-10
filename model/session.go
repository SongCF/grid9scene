package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"jhqc.com/songcf/scene/pb"
	"net"
	"time"
)

type Session struct {
	AppId   string
	Uid     int32
	ChanOut chan []byte   `json:"-"`
	Die     chan struct{} `json:"-"`

	IP          net.IP
	Conn        net.Conn `json:"-"`
	PacketCount int32    //对包进行计数
	ConnectTime time.Time
}

func (s *Session) Rsp(cmd int32, payload proto.Message) {
	if s == nil {
		return
	}
	if cmd == 0 || payload == nil {
		return
	}
	m, err := proto.Marshal(payload)
	if err != nil {
		log.Infoln("Error: Marshal payload failed!")
		return
	}
	var vsn int32 = 1
	packet := &pb.Packet{
		Cmd:     &cmd,
		Payload: m,
		Vsn:     &vsn,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		log.Infoln("Error: Marshal packet failed!")
		return
	}
	if s != nil && s.ChanOut != nil {
		log.Debugf("ack msg: %v", pb.RCode[int(cmd)])
		s.ChanOut <- data
	}
}

func (s *Session) HasLogin() bool {
	if s != nil && s.AppId != "" && s.Uid != 0 {
		return true
	}
	return false
}

func (s *Session) Clean() {
	s.AppId = ""
	s.Uid = 0
	s.PacketCount = 0
}
