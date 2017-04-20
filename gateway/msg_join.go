package gateway

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"fmt"
)

// join space
func JoinReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.JoinReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdJoinReq, pb.ErrMsgFormat)
	}

	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdJoinReq, pb.ErrNotLogin)
	}

	// get join space/x/y/angle
	joinSpaceId, x, y, angle, e := getJoinPos(s, payload)
	if e != nil {
		return pb.Error(pb.CmdJoinReq, e)
	}

	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("JoinReq(user[%v:%v]) CCPool:Get error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)

	// check already join
	oldUserInfo, e := GetUserInfo(s.AppId, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdJoinReq, e)
	}
	if joinSpaceId == oldUserInfo.SpaceId {
		return pb.Error(pb.CmdJoinReq, pb.ErrAlreadyJoinSpace)
	}

	// calc grid id
	gridWidth, gridHeight, e := GetSpaceInfo(s.AppId, joinSpaceId)
	if e != nil {
		return pb.Error(pb.CmdJoinReq, e)
	}
	joinGridId := CalcGridId(x, y, gridWidth, gridHeight)

	// modify cache data
	err = conn.Cmd("MULTI").Err
	if err != nil {
		log.Errorf("JoinReq(user[%v:%v]) Cache Cmd(MULTI) error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
	}
	//1. add to grid set
	joinSet := fmt.Sprintf(FORMAT_GRID, s.AppId, joinSpaceId, joinGridId)
	if oldUserInfo.SpaceId != NIL {
		leaveSet := fmt.Sprintf(FORMAT_GRID, s.AppId, oldUserInfo.SpaceId, oldUserInfo.GridId)
		err = conn.Cmd("SMOVE", leaveSet, joinSet, s.Uid).Err
	} else {
		err = conn.Cmd("SADD", joinSet, s.Uid).Err
	}
	if err != nil {
		log.Errorf("JoinReq(user[%v:%v]) Cache Cmd(SMOVE/SADD) error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
	}
	//2. set user hash table
	uKey := fmt.Sprintf(FORMAT_USER, s.AppId, s.Uid)
	err = conn.Cmd("HMSET", uKey, "space_id", joinSpaceId, "grid_id", joinGridId,
		"x", x, "y", y, "angle", angle, "exd", string(payload.ExData), "moveTime", 0).Err
	if err != nil {
		log.Errorf("JoinReq(user[%v:%v]) Cache Cmd(HMSET) error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
	}
	err = conn.Cmd("EXEC").Err
	if err != nil {
		log.Errorf("JoinReq(user[%v:%v]) Cache Cmd(EXEC) error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
	}

	//=====================
	//ack
	//=====================
	ack := &pb.JoinAck{}
	s.Rsp(pb.CmdJoinAck, ack)
	//=====================
	//user_list ntf
	//=====================
	roundGridIdL := RoundGridAndSelf(joinGridId)
	uidL, e := GetRoundUidList(s.AppId, joinSpaceId, roundGridIdL, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdJoinReq, e)
	}
	log.Debugf("round uid list: %v", uidL)
	for _, id := range uidL {
		conn.PipeAppend("HMGET", fmt.Sprintf(FORMAT_USER, s.AppId, id), "x", "y", "angle", "exd")
	}
	ul := make([]*pb.UserData, len(uidL))
	for i, uid := range uidL {
		resp := conn.PipeResp()
		uInfoResp, err := resp.Array()
		if err != nil || len(uInfoResp) != 4 {
			log.Errorf("JoinReq(user[%v:%v]) get userInfo array error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
		}
		tX, err0 := uInfoResp[0].Float64()
		tY, err1 := uInfoResp[1].Float64()
		tAngle, err2 := uInfoResp[2].Float64()
		tExd, err3 := uInfoResp[3].Str()
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
			log.Errorf("JoinReq(user[%v:%v]) parse userInfo error(%v,%v,%v,%v)",
				s.AppId, s.Uid, err0, err1, err2, err3)
			return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
		}
		// already except self uid by func getRoundUidList
		ttx := float32(tX)
		tty := float32(tY)
		tta := float32(tAngle)
		ul[i] = &pb.UserData{
			UserId: &uid,
			PosX: &ttx,
			PosY: &tty,
			Angle: &tta,
			ExData: []byte(tExd),
		}
	}
	ulNtf := &pb.UserListNtf{
		UserList: ul,
	}
	s.Rsp(pb.CmdUserListNtf, ulNtf)
	//=====================
	//join ntf
	//=====================
	joinNtf := &pb.JoinNtf{User:&pb.UserData{
		UserId: &s.Uid,
			PosX: &x,
			PosY: &y,
			Angle: &angle,
			ExData: payload.ExData,
		}}
	for _, uid := range uidL {
		tSe := GetSession(s.AppId, uid)
		tSe.Rsp(pb.CmdJoinNtf, joinNtf)
	}
	return 0, nil
}

func getJoinPos(s *Session, payload *pb.JoinReq) (spaceId string, x, y, angle float32, e *pb.ErrInfo) {
	spaceId = NIL
	x = DEFAULT_POS_X
	y = DEFAULT_POS_Y
	angle = DEFAULT_ANGLE
	e = nil
	if payload.GetSpaceId() != nil {
		spaceId = string(payload.GetSpaceId())
	} else {
		// query last space from db
		raw := DB.QueryRow("SELECT space_id FROM last_space WHERE app_id=? and user_id=?;",
			s.AppId, s.Uid)
		var tmpSpaceId string
		err := raw.Scan(&tmpSpaceId) // if empty, err = sql.ErrNoRows
		if err != nil {
			log.Errorf("user(%v:%v) Query db failed(select last_space)", s.AppId, s.Uid)
			e = pb.ErrQueryDBError
			return
		} else {
			spaceId = tmpSpaceId
		}
	}
	if payload.GetUseLast() {
		// query last pos from db
		raw := DB.QueryRow("SELECT x,y,angle FROM last_pos WHERE app_id=? and user_id=? and space_id=?;",
			s.AppId, s.Uid, spaceId)
		var tmpX, tmpY, tmpAngle float32
		err := raw.Scan(&tmpX, &tmpY, &tmpAngle)
		if err != nil {
			log.Errorf("user(%v:%v) Query db failed(select last_pos)", s.AppId, s.Uid)
			e = pb.ErrQueryDBError
			return
		} else {
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
	return
}

