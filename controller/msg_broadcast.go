package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
)

func BroadcastReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.BroadcastReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrMsgFormat)
	}
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrNotLogin)
	}
	// space
	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("BroadcastReq(user[%v:%v]) CCPool:Get error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdQueryPosReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)
	userInfo, e := GetUserInfo(s.AppId, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdBroadcastReq, e)
	}

	if userInfo.SpaceId != NIL && userInfo.GridId != NIL {
		roundGridIdL := RoundGridAndSelf(userInfo.GridId)
		uidL, e := GetRoundUidList(s.AppId, userInfo.SpaceId, roundGridIdL, s.Uid, conn)
		if e != nil {
			return pb.Error(pb.CmdBroadcastReq, e)
		}
		log.Debugf("round uid list: %v", uidL)
		//ack
		ack := &pb.BroadcastAck{}
		s.Rsp(pb.CmdBroadcastAck, ack)
		//ntf
		ntf := &pb.BroadcastNtf{
			UserId: &s.Uid,
			Data:   payload.GetData(),
		}
		for _, uid := range uidL {
			tSe := GetSession(s.AppId, uid)
			tSe.Rsp(pb.CmdBroadcastNtf, ntf)
		}
	} else {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrNotJoinSpace)
	}
	return 0, nil
}
