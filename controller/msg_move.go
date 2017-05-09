package controller

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
)

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

	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("MoveReq(user[%v:%v]) CCPool:Get error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)

	// space
	oldUserInfo, e := GetUserInfo(s.AppId, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdMoveReq, e)
	}
	if oldUserInfo.SpaceId == NIL || oldUserInfo.GridId == NIL {
		return pb.Error(pb.CmdMoveReq, pb.ErrNotJoinSpace)
	}
	// move time
	if payload.GetTime() < oldUserInfo.MoveTime {
		//ignore
		log.Infof("ignore move_req, req_time=%d, already_time=%d", payload.GetTime(), oldUserInfo.MoveTime)
		return 0, nil
	}
	x := payload.GetPosX()
	y := payload.GetPosY()
	angle := payload.GetAngle()

	//calc dst grid
	spaceId := oldUserInfo.SpaceId
	gridWidth, gridHeight, e := GetSpaceInfo(s.AppId, spaceId)
	if e != nil {
		return pb.Error(pb.CmdMoveReq, pb.ErrNotJoinSpace)
	}
	srcGridId := oldUserInfo.GridId
	dstGridId := CalcGridId(payload.GetPosX(), payload.GetPosY(), gridWidth, gridHeight)

	// modify cache data
	if dstGridId != srcGridId {
		srcGridSet := fmt.Sprintf(FORMAT_GRID, s.AppId, spaceId, srcGridId)
		dstGridSet := fmt.Sprintf(FORMAT_GRID, s.AppId, spaceId, dstGridId)
		err = conn.Cmd("MULTI").Err
		if err != nil {
			log.Errorf("JoinReq(user[%v:%v]) Cache Cmd(MULTI) error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdJoinReq, pb.ErrServerBusy)
		}
		//1. from src to dst grid
		err := conn.Cmd("SMOVE", srcGridSet, dstGridSet, s.Uid).Err
		if err != nil {
			log.Errorf("MoveReq(user[%v:%v]) Cache Cmd error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
		}
		//2. set user hash table
		err = conn.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, s.AppId, s.Uid),
			StrGridId, dstGridId, StrX, x, StrY, y, StrAngle, angle, StrMoveTime, payload.GetTime()).Err
		if err != nil {
			log.Errorf("MoveReq(user[%v:%v]) Cache Cmd(HMSET) error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
		}
		err = conn.Cmd("EXEC").Err
		if err != nil {
			log.Errorf("MoveReq(user[%v:%v]) Cache Cmd(EXEC) error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
		}
	} else {
		err = conn.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, s.AppId, s.Uid),
			StrX, x, StrY, y, StrAngle, angle, StrMoveTime, payload.GetTime()).Err
		if err != nil {
			log.Errorf("MoveReq(user[%v:%v]) Cache Cmd(HMSET) error(%v)", s.AppId, s.Uid, err)
			return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
		}
	}

	//ack
	ack := &pb.MoveAck{
		PosX:  payload.PosX,
		PosY:  payload.PosY,
		Angle: payload.Angle,
	}
	s.Rsp(pb.CmdMoveAck, ack)
	//move_ntf
	roundGridIdL := RoundGridAndSelf(srcGridId)
	uidL, e := GetRoundUidList(s.AppId, spaceId, roundGridIdL, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdMoveReq, e)
	}
	moveNtf := &pb.MoveNtf{
		UserId: &s.Uid,
		PosX:   &x,
		PosY:   &y,
		Angle:  &angle,
	}
	for _, uid := range uidL {
		tSe := GetSession(s.AppId, uid)
		tSe.Rsp(pb.CmdMoveNtf, moveNtf)
	}

	// cross grid  leave/join ntf   user_list ntf
	if dstGridId != srcGridId {
		oldRoundGrids := RoundGridAndSelf(srcGridId)
		newRoundGrids := RoundGridAndSelf(dstGridId)
		leaveGrids := SubGrids(oldRoundGrids, newRoundGrids)
		joinGrids := SubGrids(newRoundGrids, oldRoundGrids)
		//leave cur grid
		leaveNtfUidL, e := GetRoundUidList(s.AppId, spaceId, leaveGrids, s.Uid, conn)
		if e != nil {
			return pb.Error(pb.CmdMoveReq, e)
		}
		log.Debugf("leave ntf uid list: %v", leaveNtfUidL)
		leaveNtf := &pb.LeaveNtf{UserId: &s.Uid}
		for _, uid := range leaveNtfUidL {
			tSe := GetSession(s.AppId, uid)
			tSe.Rsp(pb.CmdLeaveNtf, leaveNtf)
		}
		//join dst grid
		joinNtfUidL, e := GetRoundUidList(s.AppId, spaceId, joinGrids, s.Uid, conn)
		if e != nil {
			return pb.Error(pb.CmdMoveReq, e)
		}
		log.Debugf("join ntf uid list: %v", joinNtfUidL)
		joinNtf := &pb.JoinNtf{User: &pb.UserData{
			UserId: &s.Uid,
			PosX:   &x,
			PosY:   &y,
			Angle:  &angle,
			ExData: oldUserInfo.ExData,
		}}
		for _, uid := range joinNtfUidL {
			tSe := GetSession(s.AppId, uid)
			tSe.Rsp(pb.CmdJoinNtf, joinNtf)
		}
		// ntf user list
		for _, id := range joinNtfUidL {
			conn.PipeAppend("HMGET", fmt.Sprintf(FORMAT_USER, s.AppId, id),
				StrX, StrY, StrAngle, StrExData)
		}
		ul := make([]*pb.UserData, len(joinNtfUidL))
		for i, uid := range joinNtfUidL {
			resp := conn.PipeResp()
			uInfoResp, err := resp.Array()
			if err != nil {
				log.Errorf("MoveReq(user[%v:%v]) get userInfo array error(%v)", s.AppId, s.Uid, err)
				return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
			}
			tX, err0 := uInfoResp[0].Float64()
			tY, err1 := uInfoResp[1].Float64()
			tAngle, err2 := uInfoResp[2].Float64()
			tExd, err3 := uInfoResp[3].Str()
			if err != nil {
				log.Errorf("MoveReq(user[%v:%v]) parse userInfo error(%v,%v,%v,%v)",
					s.AppId, s.Uid, err0, err1, err2, err3)
				return pb.Error(pb.CmdMoveReq, pb.ErrServerBusy)
			}
			// already except self uid
			ttx := float32(tX)
			tty := float32(tY)
			tta := float32(tAngle)
			ul[i] = &pb.UserData{
				UserId: &uid,
				PosX:   &ttx,
				PosY:   &tty,
				Angle:  &tta,
				ExData: []byte(tExd),
			}
		}
		ulNtf := &pb.UserListNtf{UserList: ul}
		s.Rsp(pb.CmdUserListNtf, ulNtf)
	}

	return 0, nil
}
