package gateway

import (
	"time"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	"jhqc.com/songcf/scene/pb"
	"github.com/golang/protobuf/proto"
)


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
			req, err := handleMsg(s, msg)
			s.ErrorNtf(req, err)
		case <- minTimer:
			timeWork()
			minTimer = time.After(time.Minute)
		case <- GlobalDie:
			return
		}
	}
}



func handleMsg(s *Session, m []byte) (int, *pb.ErrInfo) {
	packet := &pb.Packet{}
	err := proto.Unmarshal(m, packet)
	if err != nil {
		return 0, pb.ErrMsgFormat
	}

	if fn, ok := pb.Handlers[packet.GetCmd()]; !ok {
		return packet.GetCmd(), pb.ErrCmdNotSupport
	} else {
		return fn(s, packet.GetPayload())
	}
}