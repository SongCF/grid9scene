package controller

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	"github.com/golang/protobuf/proto"
	"jhqc.com/songcf/scene/pb"
)



func Msg2Grid(appId, spaceId, gridId string, gridMsg *GridMsg) {
	grid := GetGrid(appId, spaceId, gridId)
	if grid != nil {
		grid.PostMsg(gridMsg)
	} else {
		Msg2SpaceWait(appId, spaceId, IMSG_START_GRID, gridId)
		grid = GetGrid(appId, spaceId, gridId)
		if grid == nil {
			log.Errorf("start grid_server[%v:%v:%v] failed, ignore msg:%v\n", appId, spaceId, gridId, gridMsg)
			return
		}
		grid.PostMsg(gridMsg)
	}
}

// unsafe multi goruntine
func StartGrid(appId, spaceId, gridId string) *Grid {
	ch := make(chan struct{})
	go gridServe(appId, spaceId, gridId, ch)
	<- ch
	grid := GetGrid(appId, spaceId, gridId)
	return grid
}

func gridServe(appId, spaceId, gridId string, ch chan struct{}) {
	alreadyGrid := GetGrid(appId, spaceId, gridId)
	if alreadyGrid != nil {
		log.Infof("grid server[%v:%v:%v] already exist.", appId, spaceId, gridId)
		close(ch) //started.
		return
	}

	grid := &Grid{
		GridId: gridId,
		UidM: make(map[int32]bool),
		MsgBox: make(chan *GridMsg),
		Die: make(chan struct{}),
	}
	SetGrid(appId, spaceId, gridId, grid)
	defer func() {
		//delete cache
		if app, ok := AppL[appId]; ok {
			if space, ok := app.SpaceM[spaceId]; ok {
				delete(space.GridM, gridId)
			}
		}
		log.Infof("grid server stop. %v:%v:%v", appId, spaceId, gridId)
	}()
	close(ch) //started.

	// loop
	for {
		select {
		case m, ok := <- grid.MsgBox:
			if !ok {
				return
			}
			log.Infof("handle grid msg:%v", m)

			switch m.Cmd {
			case pb.CmdJoinReq:
				p := m.Msg.(*pb.JoinReq)
				join(appId, spaceId, gridId, m.Uid, p)
			case pb.CmdLeaveReq:
				_ = m.Msg.(*pb.LeaveReq)
				uInfo := m.ExData.(*UserInfo)
				leave(appId, spaceId, gridId, m.Uid, uInfo)
			case pb.CmdMoveReq:
				p := m.Msg.(*pb.MoveReq)
				move(appId, spaceId, gridId, m.Uid, p)
			case pb.CmdBroadcastReq:
				p := m.Msg.(*pb.BroadcastReq)
				broadcast(appId, spaceId, gridId, m.Uid, p)
			default:
				log.Warnf("unknow grid msg id = %d", m.Cmd)
			}
		case <- grid.Die:
			return
		case <- GlobalDie:
			return
		}
	}
}



func join(appId, spaceId, gridId string, uid int32, req *pb.JoinReq) {
	s := GetSession(appId, uid)
	if s == nil {
		log.Errorln("Not found session when user join")
		return
	}
	grid := GetGrid(appId, spaceId, gridId)
	if grid == nil {
		log.Errorln("Not found grid when user join")
		return
	}
	//set session
	if s.UData == nil {
		s.UData = &UserInfo{}
	}
	s.UData.SpaceId = spaceId
	s.UData.GridId = gridId
	s.UData.PosX = req.GetPosX()
	s.UData.PosY = req.GetPosY()
	s.UData.Angle = req.GetAngle()
	s.UData.ExData = req.GetExData()
	s.UData.MoveTime = 0
	//set grid uidM
	grid.UidM[uid] = true
	//ack
	ack := &pb.JoinAck{}
	s.Rsp(pb.CmdJoinAck, ack)
	//cast
	ntf := &pb.JoinNtf{
		User: &pb.UserData{
			UserId: &uid,
			PosX: &s.UData.PosX,
			PosY: &s.UData.PosY,
			Angle: &s.UData.Angle,
			ExData: s.UData.ExData,
		},
	}
	cast9grid(appId, spaceId, gridId, pb.CmdJoinNtf, ntf, uid)
}

func leave(appId, spaceId, gridId string, uid int32, uInfo *UserInfo) {
	grid := GetGrid(appId, spaceId, gridId)
	if grid == nil {
		log.Errorln("Not found grid when user leave")
		return
	}
	//set session
	s := GetSession(appId, uid)
	// session maybe closed
	//if s == nil || s.UData == nil {
	//	log.Errorln("Not found session when user leave")
	//	return
	//}
	if s != nil && s.UData != nil && s.UData.SpaceId == spaceId {
		s.UData.SpaceId = ""
		s.UData.GridId = ""
		s.UData.PosX = DEFAULT_POS_X
		s.UData.PosY = DEFAULT_POS_Y
		s.UData.Angle = DEFAULT_ANGLE
		s.UData.MoveTime = 0
	}
	//set grid uidM
	delete(grid.UidM, uid)
	//ack
	ack := &pb.LeaveAck{}
	s.Rsp(pb.CmdLeaveAck, ack)
	//cast
	ntf := &pb.LeaveNtf{UserId: &uid}
	cast9grid(appId, spaceId, gridId, pb.CmdLeaveNtf, ntf, uid)
	//db: last space pos
	tx, err := DB.Begin()
	if err != nil {
		log.Errorln("save last space and pos failed")
	} else {
		//此时 session可能已经被删除了，所以不能使用session中的位置信息
		tx.Exec("REPLACE INTO last_space(app_id,user_id,space_id) VALUES(?,?,?);",
			appId, uid, spaceId)
		tx.Exec("REPLACE INTO last_pos(app_id,user_id,space_id,x,y,angle) VALUES(?,?,?,?,?,?);",
			appId, uid, spaceId, uInfo.PosX, uInfo.PosY, uInfo.Angle)
		err = tx.Commit()
		if err != nil {
			log.Errorln("update last space and pos failed, db commit failed")
		} else {
			log.Infof("user(%v) save last space and pos success!", uid)
		}
	}
}

func move(appId, spaceId, gridId string, uid int32, req *pb.MoveReq) {
	// check move time, already check at handlers
	//
	s := GetSession(appId, uid)
	if s == nil || s.UData == nil {
		log.Errorln("Not found session when user move")
		return
	}
	oldGrid := GetGrid(appId, spaceId, gridId)
	if oldGrid == nil {
		log.Errorln("Not found grid when user move")
		return
	}
	//calc dst grid
	space := GetSpace(appId, spaceId)
	if space == nil {
		log.Errorln("not found space when user move")
		return
	}
	dstGridId := CalcGridId(req.GetPosX(), req.GetPosY(), space.GridWidth, space.GridHeight)
	dstGrid := GetGrid(appId, spaceId, dstGridId)
	if dstGrid == nil {
		Msg2SpaceWait(appId, spaceId, IMSG_START_GRID, dstGridId)
		dstGrid = GetGrid(appId, spaceId, dstGridId)
		if dstGrid == nil {
			log.Errorln("start dst grid_server[%v:%v:%v] failed.", appId, spaceId, dstGridId)
			return
		}
	}
	//set session
	s.UData.GridId = dstGridId
	s.UData.PosX = req.GetPosX()
	s.UData.PosY = req.GetPosY()
	s.UData.Angle = req.GetAngle()
	s.UData.MoveTime = req.GetTime()
	//set grid uidM
	delete(oldGrid.UidM, uid)
	dstGrid.UidM[uid] = true
	//ack
	ack := &pb.MoveAck{
		PosX: &s.UData.PosX,
		PosY: &s.UData.PosY,
		Angle: &s.UData.Angle,
	}
	s.Rsp(pb.CmdMoveAck, ack)
	//cast
	ntf := &pb.MoveNtf{
		UserId: &uid,
		PosX: &s.UData.PosX,
		PosY: &s.UData.PosY,
		Angle: &s.UData.Angle,
	}
	cast9grid(appId, spaceId, gridId, pb.CmdMoveNtf, ntf, uid)

	//cross grid move
	if dstGridId != gridId {
		oldRoundGrids := RoundGridAndSelf(gridId)
		newRoundGrids := RoundGridAndSelf(dstGridId)
		leaveGrids := SubGrids(oldRoundGrids, newRoundGrids)
		joinGrids := SubGrids(newRoundGrids, oldRoundGrids)
		//leave cur grid
		leaveNtf := &pb.LeaveNtf{UserId: &uid}
		castGrids(appId, spaceId, leaveGrids, pb.CmdLeaveNtf, leaveNtf, uid)
		//join dst grid
		joinNtf := &pb.JoinNtf{
			User: &pb.UserData{
				UserId: &uid,
				PosX: &s.UData.PosX,
				PosY: &s.UData.PosY,
				Angle: &s.UData.Angle,
				ExData: s.UData.ExData,
			},
		}
		castGrids(appId, spaceId, joinGrids, pb.CmdJoinNtf, joinNtf, uid)
		// ntf user list
		ul := []*pb.UserData{}
		for _,joinGid := range (*joinGrids) {
			joinG := GetGrid(appId, spaceId, joinGid)
			if joinG != nil {
				for joinUid := range joinG.UidM {
					joinS := GetSession(appId, joinUid)
					if joinS != nil {
						u := &pb.UserData{
							UserId: &joinUid,
							PosX: &joinS.UData.PosX,
							PosY: &joinS.UData.PosY,
							Angle: &joinS.UData.Angle,
							ExData: joinS.UData.ExData,
						}
						ul = append(ul, u)
					}
				}
			}
		}
		ulNtf := &pb.UserListNtf{UserList: ul}
		s.Rsp(pb.CmdUserListNtf, ulNtf)
	}
}

func broadcast(appId, spaceId, gridId string, uid int32, req *pb.BroadcastReq) {
	s := GetSession(appId, uid)
	if s == nil {
		log.Errorln("Not found session when user broadcast")
		return
	}
	//ack
	ack := &pb.BroadcastAck{}
	s.Rsp(pb.CmdBroadcastAck, ack)
	//cast
	ntf := &pb.BroadcastNtf{
		UserId: &uid,
		Data: req.GetData(),
	}
	cast9grid(appId, spaceId, gridId, pb.CmdBroadcastNtf, ntf, uid)
}



//=============================
//inner func
//=============================


//广播9个格子
func cast9grid(appId, spaceId, gridId string, cmd int32, ntf proto.Message, self int32) {
	gridIdL := RoundGridAndSelf(gridId)
	castGrids(appId, spaceId, gridIdL, cmd, ntf, self)
}

func castGrids(appId, spaceId string, gridIdL *[]string, cmd int32, ntf proto.Message, self int32) {
	for i := 0; i < len(*gridIdL); i ++ {
		grid := GetGrid(appId, spaceId, (*gridIdL)[i])
		if grid != nil {
			for uid := range grid.UidM {
				if uid != self {
					s := GetSession(appId, uid)
					if s != nil {
						s.Rsp(cmd, ntf)
					}
				}
			}
		}
	}
}
