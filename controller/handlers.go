package controller

import (
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
)

//define this func by yourself
var Handlers = map[int32]func(*Session, []byte) (int32, proto.Message){
	10020: LoginReq,     // 用于登录服务器 验证
	20010: JoinReq,      // 加入场景请求
	20020: LeaveReq,     // 离开场景请求
	20030: MoveReq,      // 场景移动请求
	20040: BroadcastReq, // 场景广播请求
	20050: QueryPosReq,  // 查询位置请求
	20061: HeartbeatReq, // 系统预留，用于表示心跳
}
