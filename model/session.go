package model

import (
	"time"
	"net"
	"jhqc.com/songcf/scene/pb"
	"github.com/golang/protobuf/proto"
	log "github.com/Sirupsen/logrus"
)

type Session struct {
	AppId       string
	Uid         int32
	ChanOut     chan []byte `json:"-"`
	Die 	    chan struct{} `json:"-"`

	IP          net.IP
	Conn        net.Conn `json:"-"`
	PacketCount int32         //对包进行计数
	ConnectTime time.Time

	UData       *UserInfo
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
		Cmd: &cmd,
		Payload: m,
		Vsn: &vsn,
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
	if s != nil && s.Die != nil {
		close(s.Die)
		s.Die = nil
		if app, ok := AppL[s.AppId]; ok {
			delete(app.SessionM, s.Uid)
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
	if app, ok := AppL[appId]; ok {
		app.SessionM[uid] = s
	} else {
		log.Errorln("not found app")
	}
}

func GetSession(appId string, uid int32) *Session {
	if app, ok := AppL[appId]; ok {
		if s, ok := app.SessionM[uid]; ok {
			return s
		}
	}
	return nil
}

//get user current space_id, grid x y
func GetUserData(appId string, uid int32) *UserInfo {
	if app, ok := AppL[appId]; ok {
		if s, ok := app.SessionM[uid]; ok {
			return s.UData
		}
	}
	return nil
}



