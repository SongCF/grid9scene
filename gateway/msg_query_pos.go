package gateway

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
)

func QueryPosReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.QueryPosReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdQueryPosReq, pb.ErrMsgFormat)
	}
	if payload.UserId == nil {
		return pb.Error(pb.CmdQueryPosReq, pb.ErrMissParam)
	}

	//query cache
	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("QueryPosReq(user[%v:%v]) CCPool:Get error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdQueryPosReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)
	userInfo, e := GetUserInfo(s.AppId, payload.GetUserId(), conn)
	if e != nil {
		return pb.Error(pb.CmdQueryPosReq, e)
	}
	ack := &pb.QueryPosAck{
		SpaceId: []byte(userInfo.SpaceId),
		PosX:    &userInfo.PosX,
		PosY:    &userInfo.PosY,
		Angle:   &userInfo.Angle,
	}
	return pb.CmdQueryPosAck, ack
}
