package gateway

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "github.com/SongCF/scene/controller"
	. "github.com/SongCF/scene/model"
	"github.com/SongCF/scene/pb"
	. "github.com/SongCF/scene/util"
	"time"
)

func agent(s *Session, in chan []byte) {
	defer log.Debug("---session agent end.")
	defer RecoverPanic()
	defer CloseSession(s)

	minTimer := time.After(time.Minute)
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}
			s.PacketCount++
			cmd, payload := handleMsg(s, msg)
			s.Rsp(cmd, payload)
		case <-minTimer:
			timeWork()
			minTimer = time.After(time.Minute)
		case <-s.Die:
			return
		case <-GlobalDie:
			return
		}
	}
}

func handleMsg(s *Session, m []byte) (int32, proto.Message) {
	defer RecoverPanic()
	if s == nil {
		return pb.Error(0, pb.ErrUser)
	}
	packet := &pb.Packet{}
	err := proto.Unmarshal(m, packet)
	if err != nil {
		return pb.Error(0, pb.ErrMsgFormat)
	}
	log.Debugf("(%v) req msg: %v", s.Uid, pb.RCode[int(packet.GetCmd())])
	if fn, ok := Handlers[packet.GetCmd()]; ok {
		return fn(s, packet.GetPayload())
	} else {
		return pb.Error(packet.GetCmd(), pb.ErrCmdNotSupport)
	}
}

func timeWork() {
	//TODO something
	//log.Info("on minute timer")
}
