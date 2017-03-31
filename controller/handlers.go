package controller

import (
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	log "github.com/Sirupsen/logrus"
	"jhqc.com/songcf/scene/pb"
)


//define this func by yourself
var Handlers = map[int32]func(*Session, []byte) (int32, proto.Message) {
	10020:   LoginReq,                // 用于登录服务器 验证
	20010:   JoinReq,                 // 加入场景请求
	20020:   LeaveReq,                // 离开场景请求
	20030:   MoveReq,                 // 场景移动请求
	20040:   BroadcastReq,            // 场景广播请求
	20050:   QueryPosReq,             // 查询位置请求
	20061:   HeartbeatReq,            // 系统预留，用于表示心跳
}


//user login
func LoginReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.LoginReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdLoginReq, pb.ErrMsgFormat)
	}
	// check the session already login
	if s.HasLogin() {
		return pb.Error(pb.CmdLoginReq, pb.ErrDuplicateLogin)
	}
	// 1.check has app
	if !HasApp(string(payload.AppId)) {
		return pb.Error(pb.CmdLoginReq, pb.ErrAppNotExist)
	}
	// kick out other session login by the uid if had
	oldSess := GetSession(string(payload.AppId), payload.GetUserId())
	if oldSess != nil {
		t := OFFLINE_TYPE_OTHER_LOGIN
		offlineMsg := & pb.OfflineNtf{
			Type: &t,
		}
		oldSess.Rsp(pb.CmdOfflineNtf, offlineMsg)
		return 0, nil
	}
	// set session user info
	s.AppId = string(payload.AppId)
	s.Uid = payload.GetUserId()
	SetSession(string(payload.AppId), payload.GetUserId(), s)
	rsp := &pb.LoginAck{}
	return pb.CmdLoginAck, rsp
}

// join space
func JoinReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.JoinReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdJoinReq, pb.ErrMsgFormat)
	}
	spaceId := string(payload.SpaceId)
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdJoinReq, pb.ErrNotLogin)
	}
	// check already join
	oldUserdata := GetUserData(s.AppId, s.Uid)
	if oldUserdata != nil && spaceId == oldUserdata.SpaceId {
		return pb.Error(pb.CmdJoinReq, pb.ErrAlreadyJoinSpace)
	}
	// get x, y, angle
	x := DEFAULT_POS_X
	y := DEFAULT_POS_Y
	angle := DEFAULT_ANGLE
	if payload.Angle != nil {
		angle = payload.GetAngle()
	}
	if payload.GetUseLast() {
		// query last space from db
		raw := DB.QueryRow("SELECT space_id FROM last_space WHERE app_id=? and user_id=?;",
			s.AppId, s.Uid)
		var tmpSpaceId string
		err = raw.Scan(&tmpSpaceId) // if empty, err = sql.ErrNoRows
		if err == nil {
			spaceId = tmpSpaceId
		}
		// query last pos from db
		raw = DB.QueryRow("SELECT x,y,angle FROM last_pos WHERE app_id=? and user_id=? and space_id=?;",
			s.AppId, s.Uid, payload.GetSpaceId())
		var tmpX, tmpY, tmpAngle float32
		err = raw.Scan(&tmpX, &tmpY, &tmpAngle)
		if err == nil {
			x = tmpX
			y = tmpY
			angle = tmpAngle
		}
	} else {
		if payload.PosX != nil {
			x = payload.GetPosX()
		}
		if payload.PosY != nil {
			y = payload.GetPosY()
		}
		if payload.Angle != nil {
			angle = payload.GetAngle()
		}
	}
	// calc grid id
	spaceInfo := GetSpace(s.AppId, spaceId)
	if spaceInfo == nil {
		return pb.Error(pb.CmdJoinReq, pb.ErrSpaceNotExist)
	}
	joinGridId := CalcGridId(x, y, spaceInfo.GridWidth, spaceInfo.GridHeight)
	// leave old space before join new space
	if oldUserdata.SpaceId != "" {
		//leave
		leaveReq := &pb.LeaveReq{SpaceId: []byte(oldUserdata.SpaceId)}
		gridMsg := &GridMsg{
			Uid: s.Uid,
			Cmd: pb.CmdLeaveReq,
			Msg: leaveReq,
		}
		Msg2Grid(s.AppId, oldUserdata.SpaceId, oldUserdata.GridId, gridMsg)
	}
	// join
	payload.SpaceId = []byte(spaceId)
	payload.PosX = &x
	payload.PosY = &y
	payload.Angle = &angle
	gridMsg := &GridMsg{
		Uid: s.Uid,
		Cmd: pb.CmdJoinReq,
		Msg: payload,
	}
	Msg2Grid(s.AppId, oldUserdata.SpaceId, joinGridId, gridMsg)
	return 0, nil
}

// leave space
func LeaveReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.LeaveReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdLeaveReq, pb.ErrMsgFormat)
	}
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdJoinReq, pb.ErrNotLogin)
	}
	// has joined?
	if s.UData != nil && s.UData.SpaceId != "" {
		if s.UData.SpaceId == string(payload.GetSpaceId()) ||
			payload.SpaceId == nil {
			//leave request space_id
			//leave current space_id
			payload.SpaceId = []byte(s.UData.SpaceId)
			gridMsg := &GridMsg{
				Uid: s.Uid,
				Cmd: pb.CmdLeaveReq,
				Msg: payload,
			}
			Msg2Grid(s.AppId, s.UData.SpaceId, s.UData.GridId, gridMsg)
			//db: last space pos
			tx, err := DB.Begin()
			if err != nil {
				log.Errorln("save last space and pos failed")
			} else {
				tx.Exec("REPLACE INTO last_space(app_id,user_id,space_id) VALUES(?,?,?);",
					s.AppId, s.Uid, s.UData.SpaceId)
				tx.Exec("REPLACE INTO last_pos(app_id,user_id,space_id,x,y,angle) VALUES(?,?,?,?,?,?);",
					s.AppId, s.Uid, s.UData.SpaceId, s.UData.PosX, s.UData.PosY, s.UData.Angle)
				err = tx.Commit()
				if err != nil {
					log.Errorln("delete app failed, db commit failed")
				}
			}
		}
	}
	//ignore
	return 0, nil
}

func MoveReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.MoveReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdMoveReq, pb.ErrMsgFormat)
	}
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdMoveReq, pb.ErrNotLogin)
	}
	// space
	if s.UData == nil || s.UData.SpaceId == "" || s.UData.GridId == "" {
		return pb.Error(pb.CmdMoveReq, pb.ErrNotJoinSpace)
	}
	// move time
	if payload.GetTime() < s.UData.MoveTime {
		//ignore
		log.Infof("ignore move_req, req_time=%d, already_time=%d\n", payload.GetTime(), s.UData.MoveTime)
		return 0, nil
	}
	gridMsg := &GridMsg{
		Uid: s.Uid,
		Cmd: pb.CmdMoveReq,
		Msg: payload,
	}
	Msg2Grid(s.AppId, s.UData.SpaceId, s.UData.GridId, gridMsg)
	return 0, nil
}

func BroadcastReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.MoveReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrMsgFormat)
	}
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrNotLogin)
	}
	// space
	if s.UData != nil && s.UData.SpaceId != "" && s.UData.GridId != "" {
		gridMsg := &GridMsg{
			Uid: s.Uid,
			Cmd: pb.CmdBroadcastReq,
			Msg: payload,
		}
		Msg2Grid(s.AppId, s.UData.SpaceId, s.UData.GridId, gridMsg)
		return 0, nil
	} else {
		return pb.Error(pb.CmdBroadcastReq, pb.ErrNotJoinSpace)
	}
}

func QueryPosReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.QueryPosReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdQueryPosReq, pb.ErrMsgFormat)
	}
	if payload.UserId == nil {
		return pb.Error(pb.CmdQueryPosReq, pb.ErrMissParam)
	}
	tarSess := GetSession(s.AppId, payload.GetUserId())
	if tarSess == nil || tarSess.UData == nil {
		return pb.Error(pb.CmdQueryPosReq, pb.ErrUserOffline)
	}
	ack := &pb.QueryPosAck{
		SpaceId: []byte(tarSess.UData.SpaceId),
		PosX: &tarSess.UData.PosX,
		PosY: &tarSess.UData.PosY,
		Angle: &tarSess.UData.Angle,
	}
	return pb.CmdQueryPosAck, ack
}

func HeartbeatReq(_ *Session, m []byte) (int32, proto.Message) {
	payload := &pb.HeartbeatReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdHeartbeatReq, pb.ErrMsgFormat)
	}
	rsp := &pb.HeartbeatAck{}
	return pb.CmdHeartbeatAck, rsp
}