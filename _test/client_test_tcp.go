package _test

import (
	"fmt"
	"jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"testing"
	"time"
)

func TestAllPB(httpAddr, tcpAddr string, t *testing.T) {
	initAppSpace(httpAddr)

	time.Sleep(time.Second * 1)
	tHeartbeat(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tLogin(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tJoin(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tLeave(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tMove(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tBroadcast(tcpAddr, t)

	time.Sleep(time.Second * 1)
	tQueryPos(tcpAddr, t)

	//end:  waiting for print debug log.
	time.Sleep(time.Second * 1)
}

func tHeartbeat(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("HEARTBEAT")
	fmt.Println("============================================")
	wCh, rCh := startClient(tcpAddr)
	defer closeClient(wCh)
	wCh <- heartbeat()
	rsp := getMsg(rCh)
	unpack(pb.CmdHeartbeatAck, rsp)
}

func tLogin(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("LOGIN")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr)
	defer closeClient(wCh2)
	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp := getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
	//重复登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrDuplicateLogin.Id, t)
	//顶号登录
	wCh2 <- login(T_APP_ID, T_USER_ID)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdLoginAck, c2Rsp)
	c1Rsp = getMsg(rCh1)
	offlineNtf := unpack(pb.CmdOfflineNtf, c1Rsp).(*pb.OfflineNtf)
	assert(offlineNtf.GetType() == model.OFFLINE_TYPE_OTHER_LOGIN, t)
	//应用不存在
	wCh1 <- login("not_exist_app_id", T_USER_ID2)
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAppNotExist.Id, t)
	//保持c1c2正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID2)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
}

func tJoin(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("JOIN")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr)
	defer closeClient(wCh2)

	var x, y float32 = 111, 111

	//未登陆时join
	wCh1 <- join(T_SPACE_ID, x, y)
	c1Rsp := getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
	wCh2 <- login(T_APP_ID, T_USER_ID2)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdLoginAck, c2Rsp)

	//space不存在
	wCh1 <- join("not_exist_space_id", x, y)
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrSpaceNotExist.Id, t)

	//join
	wCh1 <- join(T_SPACE_ID, x, y)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdJoinAck, c1Rsp)
	//user list
	c1Rsp = getMsg(rCh1)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)

	//join last
	wCh1 <- leaveOpt()
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLeaveAck, c1Rsp)
	wCh1 <- joinLastPos(T_SPACE_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdJoinAck, c1Rsp)
	//user list
	c1Rsp = getMsg(rCh1)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)

	//重复join
	wCh1 <- join(T_SPACE_ID, x, y)
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAlreadyJoinSpace.Id, t)

	//join2
	wCh2 <- join(T_SPACE_ID, x, y)
	c2Rsp = getMsg(rCh2)
	unpack(pb.CmdJoinAck, c2Rsp)
	//user list
	c2Rsp = getMsg(rCh2)
	ulNtf = unpack(pb.CmdUserListNtf, c2Rsp).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID, t)
	//join_ntf
	c1Rsp = getMsg(rCh1)
	joinNtf := unpack(pb.CmdJoinNtf, c1Rsp).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID2, t)
}

func tLeave(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("LEAVE")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr)
	defer closeClient(wCh2)

	//未登陆时leave
	wCh1 <- leave(T_SPACE_ID)
	c1Rsp := getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)
	wCh1 <- leaveOpt()
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
	wCh2 <- login(T_APP_ID, T_USER_ID2)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdLoginAck, c2Rsp)

	//未join时leave
	wCh1 <- leave(T_SPACE_ID) //不应该有msg包返回
	wCh1 <- leaveOpt()          //不应该有msg包返回

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//leave 未加入的space
	wCh1 <- leave("not_joined_space_id")
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//leave
	wCh1 <- leave(T_SPACE_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLeaveAck, c1Rsp)
	//leave_ntf
	c2Rsp = getMsg(rCh2)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)

	//leave opt
	wCh2 <- leaveOpt()
	c2Rsp = getMsg(rCh2)
	unpack(pb.CmdLeaveAck, c2Rsp)
}

func tMove(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("MOVE")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr)
	defer closeClient(wCh2)

	//未登录时move
	wCh1 <- move(0, 1, 1)
	c1Rsp := getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
	wCh2 <- login(T_APP_ID, T_USER_ID2)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdLoginAck, c2Rsp)

	var x1, y1 float32 = 1, 2
	var x2, y2 float32 = (x1 + T_GRID_W), (y1 + T_GRID_H)
	var x3, y3 float32 = (x2 + T_GRID_W), (y2 + T_GRID_H)
	var time int32 = 0

	//未join时move
	wCh1 <- move(time, x1, y1)
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//move
	time = 10
	wCh1 <- move(time, x1, y1)
	c1Rsp = getMsg(rCh1)
	moveAck := unpack(pb.CmdMoveAck, c1Rsp).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x1, t)
	assert(moveAck.GetPosY() == y1, t)
	//move_ntf
	c2Rsp = getMsg(rCh2)
	moveNtf := unpack(pb.CmdMoveNtf, c2Rsp).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x1, t)
	assert(moveNtf.GetPosY() == y1, t)

	//time ignore
	time = 9
	wCh1 <- move(time, x3, y3)

	//across move
	time = 11
	wCh1 <- move(time, x2, y2)
	c1Rsp = getMsg(rCh1)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x2, t)
	assert(moveAck.GetPosY() == y2, t)
	c1Rsp = getMsg(rCh1)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//move_ntf
	c2Rsp = getMsg(rCh2)
	moveNtf = unpack(pb.CmdMoveNtf, c2Rsp).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x2, t)
	assert(moveNtf.GetPosY() == y2, t)

	//across move again
	time = 12
	wCh1 <- move(time, x3, y3)
	c1Rsp = getMsg(rCh1)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x3, t)
	assert(moveAck.GetPosY() == y3, t)
	c1Rsp = getMsg(rCh1)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//move_ntf
	c2Rsp = getMsg(rCh2)
	moveNtf = unpack(pb.CmdMoveNtf, c2Rsp).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x3, t)
	assert(moveNtf.GetPosY() == y3, t)
	//leave_ntf (离开了九宫格范围)
	c2Rsp = getMsg(rCh2)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)

	//move back
	time = 13
	wCh1 <- move(time, x2, y2)
	c1Rsp = getMsg(rCh1)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x2, t)
	assert(moveAck.GetPosY() == y2, t)
	c1Rsp = getMsg(rCh1)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID2, t)
	//join_ntf (进入了九宫格范围)
	c2Rsp = getMsg(rCh2)
	joinNtf := unpack(pb.CmdJoinNtf, c2Rsp).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID, t)
	assert(joinNtf.GetUser().GetPosX() == x2, t)
	assert(joinNtf.GetUser().GetPosY() == y2, t)
}

func tBroadcast(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("BROADCAST")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr)
	defer closeClient(wCh2)

	//未登录
	wCh1 <- broadcast()
	c1Rsp := getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)
	wCh2 <- login(T_APP_ID, T_USER_ID2)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdLoginAck, c2Rsp)

	//未join
	wCh1 <- broadcast()
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//广播1
	wCh1 <- broadcast()
	c1Rsp = getMsg(rCh1)
	_ = unpack(pb.CmdBroadcastAck, c1Rsp).(*pb.BroadcastAck)
	//broadcast_ntf
	c2Rsp = getMsg(rCh2)
	bcNtf := unpack(pb.CmdBroadcastNtf, c2Rsp).(*pb.BroadcastNtf)
	assert(bcNtf.GetUserId() == T_USER_ID, t)

	//广播2
	//leave
	wCh1 <- leave(T_SPACE_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLeaveAck, c1Rsp)
	//leave_ntf
	c2Rsp = getMsg(rCh2)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)
	//重新join
	//join
	wCh1 <- join(T_SPACE_ID, T_GRID_W*10, T_GRID_H*10)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdJoinAck, c1Rsp)
	//user list
	c1Rsp = getMsg(rCh1)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//广播
	wCh1 <- broadcast()
	c1Rsp = getMsg(rCh1)
	_ = unpack(pb.CmdBroadcastAck, c1Rsp).(*pb.BroadcastAck)
	//none broadcast_ntf
	getNilMsg(rCh2)
}

func tQueryPos(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("QUERY_POS")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr)
	defer closeClient(wCh1)

	//自己未登录
	wCh1 <- querypos(T_USER_ID)
	c1Rsp := getMsg(rCh1)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//login
	wCh1 <- login(T_APP_ID, T_USER_ID)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdLoginAck, c1Rsp)

	//不在线
	wCh1 <- querypos(T_USER_ID2)
	c1Rsp = getMsg(rCh1)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrUserOffline.Id, t)

	//query
	wCh1 <- querypos(T_USER_ID)
	c1Rsp = getMsg(rCh1)
	qAck := unpack(pb.CmdQueryPosAck, c1Rsp).(*pb.QueryPosAck)
	assert(string(qAck.GetSpaceId()) == model.NIL, t)
	assert(qAck.GetPosX() == model.DEFAULT_POS_X, t)
	assert(qAck.GetPosY() == model.DEFAULT_POS_Y, t)
	assert(qAck.GetAngle() == model.DEFAULT_ANGLE, t)

	//join
	var x, y float32 = 190, 910
	wCh1 <- join(T_SPACE_ID, x, y)
	c1Rsp = getMsg(rCh1)
	unpack(pb.CmdJoinAck, c1Rsp)
	//user list
	c1Rsp = getMsg(rCh1)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//query
	wCh1 <- querypos(T_USER_ID)
	c1Rsp = getMsg(rCh1)
	qAck = unpack(pb.CmdQueryPosAck, c1Rsp).(*pb.QueryPosAck)
	assert(string(qAck.GetSpaceId()) == T_SPACE_ID, t)
	assert(qAck.GetPosX() == x, t)
	assert(qAck.GetPosY() == y, t)
	assert(qAck.GetAngle() == model.DEFAULT_ANGLE, t)
}

//两个user join
func joinSpace(rCh1, wCh1, rCh2, wCh2 chan []byte, t *testing.T) {
	//join
	wCh1 <- join(T_SPACE_ID, 0, 0)
	c1Rsp := getMsg(rCh1)
	unpack(pb.CmdJoinAck, c1Rsp)
	//user list
	c1Rsp = getMsg(rCh1)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//join2
	wCh2 <- join(T_SPACE_ID, 0, 0)
	c2Rsp := getMsg(rCh2)
	unpack(pb.CmdJoinAck, c2Rsp)
	//user list
	c2Rsp = getMsg(rCh2)
	ulNtf = unpack(pb.CmdUserListNtf, c2Rsp).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID, t)
	//join_ntf
	c1Rsp = getMsg(rCh1)
	joinNtf := unpack(pb.CmdJoinNtf, c1Rsp).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID2, t)
}

