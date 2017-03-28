package gateway

import (
	"time"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/controller"
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
			cmd, payload := handleMsg(s, msg)
			s.Rsp(cmd, payload)
		case <- minTimer:
			timeWork()
			minTimer = time.After(time.Minute)
		case <- GlobalDie:
			return
		}
	}
}



func handleMsg(s *Session, m []byte) (int32, proto.Message) {
	if s == nil {
		//TODO tick out
		return pb.Error(0, pb.ErrUser)
	}
	packet := &pb.Packet{}
	err := proto.Unmarshal(m, packet)
	if err != nil {
		return pb.Error(0, pb.ErrMsgFormat)
	}
	if fn, ok := Handlers[packet.GetCmd()]; ok {
		return fn(s, packet.GetPayload())
	} else {
		return pb.Error(packet.GetCmd(), pb.ErrCmdNotSupport)
	}
}

func timeWork() {
	//TODO something
	log.Info("on minute timer")
}
