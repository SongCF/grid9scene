package gateway

import (
	"time"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
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
			msgHandler(s.AppId, s.Uid, msg)
		case <- minTimer:
			timeWork()
			minTimer = time.After(time.Minute)
		case <- GlobalDie:
			return
		}
	}
}



func msgHandler(appId string, uid int32, m []byte) {
	cmd := 1
	switch cmd {
	case 1:

	case 2:
		//s := AppList[appId].sessionList[uid]
		//msg2Grid(s.appId, s.spaceId, s.gridId, m)
	default:
	}
}