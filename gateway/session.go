package gateway

import (
	"time"
	"net"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"github.com/golang/protobuf/proto"
	log "github.com/Sirupsen/logrus"
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

func (s *Session) Rsp(cmd int, payload []byte) {
	if s == nil {
		return
	}
	packet := pb.Packet{
		Cmd: cmd,
		Payload: payload,
		Vsn: 1,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		log.Infoln("Error: Marshal packet failed!")
		return
	}
	s.ChanOut <- data
}

func (s *Session) ErrorNtf(req int, errInfo *pb.ErrInfo) {
	if s == nil {
		return
	}
	e := &pb.ErrorNtf{
		Code: errInfo.Id,
		Msg: errInfo.Desc,
		Req: req,
	}
	payload, err := proto.Marshal(e)
	if err != nil {
		log.Errorf("marshal errorntf failed: %v, errInfo:%v", err, errInfo)
		return
	}

}




// mem cache
type AppInfo struct {
	SpaceM map[string]*Space      // spaceId : Space
	SessionM map[int32]*Session // uid : Session
}
var (
	AppInfoL = make(map[string]*AppInfo)
)
