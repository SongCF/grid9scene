package gateway

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"fmt"
)


// leave space
func LeaveReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.LeaveReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdLeaveReq, pb.ErrMsgFormat)
	}
	return Leave(s, payload)
}

func Leave(s *Session, payload *pb.LeaveReq) (int32, proto.Message) {
	// check login
	if !s.HasLogin() {
		return pb.Error(pb.CmdLeaveReq, pb.ErrNotLogin)
	}

	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("LeaveReq(user[%v:%v]) CCPool:Get error(%v)", s.AppId, s.Uid, err)
		return pb.Error(pb.CmdLeaveReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)

	// has joined?
	userInfo, e := GetUserInfo(s.AppId, s.Uid, conn)
	if e != nil {
		return pb.Error(pb.CmdLeaveReq, e)
	}

	if userInfo.SpaceId != NIL && userInfo.GridId != NIL {
		if payload.SpaceId == nil || userInfo.SpaceId == string(payload.GetSpaceId()) {
			// modify cache data
			err = conn.Cmd("MULTI").Err
			if err != nil {
				log.Errorf("LeaveReq(user[%v:%v]) Cache Cmd(MULTI) error(%v)", s.AppId, s.Uid, err)
				return pb.Error(pb.CmdLeaveReq, pb.ErrServerBusy)
			}
			//1. leave from grid set
			leaveSet := fmt.Sprintf(FORMAT_GRID, s.AppId, userInfo.SpaceId, userInfo.GridId)
			err = conn.Cmd("SREM", leaveSet, s.Uid).Err
			if err != nil {
				log.Errorf("LeaveReq(user[%v:%v]) Cache Cmd(SREM) error(%v)", s.AppId, s.Uid, err)
				return pb.Error(pb.CmdLeaveReq, pb.ErrServerBusy)
			}
			//2. set user hash table
			e := ResetUserInfo(s.AppId, s.Uid, conn)
			if e != nil {
				log.Errorf("LeaveReq(user[%v:%v]) ResetUserInfo error(%v)", s.AppId, s.Uid, err)
				return pb.Error(pb.CmdLeaveReq, pb.ErrServerBusy)
			}
			err = conn.Cmd("EXEC").Err
			if err != nil {
				log.Errorf("LeaveReq(user[%v:%v]) Cache Cmd(EXEC) error(%v)", s.AppId, s.Uid, err)
				return pb.Error(pb.CmdLeaveReq, pb.ErrServerBusy)
			}

			//db: last space pos
			tx, err := DB.Begin()
			if err != nil {
				log.Errorln("save last space and pos failed")
			} else {
				tx.Exec("REPLACE INTO last_space(app_id,user_id,space_id) VALUES(?,?,?);",
					s.AppId, s.Uid, userInfo.SpaceId)
				tx.Exec("REPLACE INTO last_pos(app_id,user_id,space_id,x,y,angle) VALUES(?,?,?,?,?,?);",
					s.AppId, s.Uid, userInfo.SpaceId, userInfo.PosX, userInfo.PosY, userInfo.Angle)
				err = tx.Commit()
				if err != nil {
					log.Errorln("update last space and pos failed, db commit failed")
				} else {
					log.Infof("user(%v) save last space and pos success!", s.Uid)
				}
			}

			//ack
			ack := &pb.LeaveAck{}
			s.Rsp(pb.CmdLeaveAck, ack)
			//ntf
			roundGridIdL := RoundGridAndSelf(userInfo.GridId)
			uidL, e := GetRoundUidList(s.AppId, userInfo.SpaceId, roundGridIdL, s.Uid, conn)
			if e != nil {
				return pb.Error(pb.CmdLeaveReq, e)
			}
			log.Debugf("round uid list: %v", uidL)
			leaveNtf := &pb.LeaveNtf{UserId: &s.Uid}
			for _, uid := range uidL {
				tSe := GetSession(s.AppId, uid)
				tSe.Rsp(pb.CmdLeaveNtf, leaveNtf)
			}
		} else {
			return pb.Error(pb.CmdLeaveReq, pb.ErrNotJoinSpace)
		}
	}
	return 0, nil
}

