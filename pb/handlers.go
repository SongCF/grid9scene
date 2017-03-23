package pb

import (
	. "jhqc.com/songcf/scene/types"
)


var Handlers = map[int]func(*Session, []byte){
	10020: Login, //"login_req",     // 用于登录服务器 验证
	20010: Join, //"join_req",      // 加入场景请求
	20020: Leave, //"leave_req",     // 离开场景请求
	20030: Move, //"move_req",      // 场景移动请求
	20040: Broadcast, //"broadcast_req", // 场景广播请求
	20050: QueryPos, //"query_pos_req", // 查询位置请求
	20061: Heartbeat, //"heartbeat_req", // 系统预留，用于表示心跳
}


func Login(s *Session, msg []byte) {

}

func Join(s *Session, msg []byte) {

}

func Leave(s *Session, msg []byte) {

}

func Move(s *Session, msg []byte) {

}

func Broadcast(s *Session, msg []byte) {

}

func QueryPos(s *Session, msg []byte) {

}

func Heartbeat(s *Session, msg []byte) {

}