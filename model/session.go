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
	ChanOut     chan []byte

	IP          net.IP
	Conn        net.Conn
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
	packet := &pb.Packet{
		Cmd: &cmd,
		Payload: m,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		log.Infoln("Error: Marshal packet failed!")
		return
	}
	s.ChanOut <- data
}

func (s *Session) HasLogin() bool {
	if s != nil && s.AppId != "" && s.Uid != 0 {
		return true
	}
	return false
}


func SetSession(appId string, uid int32, s *Session) {
	AppInfoL[appId].SessionM[uid] = s
}

func GetSession(appId string, uid int32) *Session {
	s, ok := AppInfoL[appId].SessionM[uid]
	if ok {
		return s
	}
	return nil
}

func GetSpace(appId, spaceId string) *Space {
	s, ok := AppInfoL[appId].SpaceM[spaceId]
	if ok {
		return s
	}
	return nil
}

//get user current space_id, grid x y
func GetUserData(appId string, uid int32) *UserInfo {
	s, ok := AppInfoL[appId].SessionM[uid]
	if ok {
		return s.UData
	}
	return nil
}



