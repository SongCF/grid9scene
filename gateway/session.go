package gateway

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



var (
	SessionPool = map[string]map[int32]*Session{}
)



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
		s.ChanOut <- data
	}
}

func (s *Session) Close() {
	if s != nil {
		select {
		case <-s.Die: //check already closed
		default:
			close(s.Die)
			// post leave
			Leave(s, &pb.LeaveReq{})
			// delete session
			if app, ok := SessionPool[s.AppId]; ok {
				delete(app, s.Uid)
			}
		}
	}
}

func (s *Session) HasLogin() bool {
	if s != nil && s.AppId != "" && s.Uid != 0 {
		return true
	}
	return false
}

func SetSession(appId string, uid int32, s *Session) {
	app, ok := SessionPool[appId]
	if !ok {
		SessionPool[appId] = map[int32]*Session{}
		app = SessionPool[appId]
	}
	app[uid] = s
}

func GetSession(appId string, uid int32) *Session {
	if app, ok := SessionPool[appId]; ok {
		if s, ok := app[uid]; ok {
			return s
		}
	}
	return nil
}
