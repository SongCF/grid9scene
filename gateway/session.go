package gateway

import (
	"time"
	"net"
	. "jhqc.com/songcf/scene/model"
)

type Session struct {
	AppId	     string
	Uid          int32
	ChanOut      chan []byte

	IP           net.IP
	Conn         net.Conn
	PacketCount  int32         //对包进行计数
	ConnectTime  time.Time

	UserData     *UserData
}


// mem cache
type AppInfo struct {
	SpaceM map[string]*Space      // spaceId : Space
	SessionM map[int32]*Session // uid : Session
}
var (
	AppInfoL = make(map[string]*AppInfo)
)
