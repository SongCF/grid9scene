package controller

import (
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
)

func HeartbeatReq(_ *Session, m []byte) (int32, proto.Message) {
	payload := &pb.HeartbeatReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdHeartbeatReq, pb.ErrMsgFormat)
	}
	rsp := &pb.HeartbeatAck{}
	return pb.CmdHeartbeatAck, rsp
}
