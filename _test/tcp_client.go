package _test

import (
	"encoding/binary"
	"fmt"
	"io"
	"jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"net"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func testAllMsg(httpAddr, tcpAddr string, t *testing.T) {
	initAppSpace(httpAddr, t)

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
	wCh, rCh := startClient(tcpAddr, t)
	defer closeClient(wCh)
	wCh <- heartbeat(t)
	rsp := getMsg(rCh, t)
	unpack(pb.CmdHeartbeatAck, rsp, t)
}

func tLogin(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("LOGIN")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr, t)
	defer closeClient(wCh2)
	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp := getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
	//重复登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrDuplicateLogin.Id, t)
	//顶号登录
	wCh2 <- login(T_APP_ID, T_USER_ID, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)
	c1Rsp = getMsg(rCh1, t)
	offlineNtf := unpack(pb.CmdOfflineNtf, c1Rsp, t).(*pb.OfflineNtf)
	assert(offlineNtf.GetType() == model.OFFLINE_TYPE_OTHER_LOGIN, t)
	//应用不存在
	wCh1 <- login("not_exist_app_id", T_USER_ID2, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAppNotExist.Id, t)
	//保持c1c2正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID2, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
}

func tJoin(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("JOIN")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr, t)
	defer closeClient(wCh2)

	var x, y float32 = 111, 111

	//未登陆时join
	wCh1 <- join(T_SPACE_ID, x, y, t)
	c1Rsp := getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
	wCh2 <- login(T_APP_ID, T_USER_ID2, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)

	//space不存在
	wCh1 <- join("not_exist_space_id", x, y, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrSpaceNotExist.Id, t)

	//join
	wCh1 <- join(T_SPACE_ID, x, y, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)

	//join last
	wCh1 <- leaveOpt(t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLeaveAck, c1Rsp, t)
	wCh1 <- joinLastPos(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)

	//重复join
	wCh1 <- join(T_SPACE_ID, x, y, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAlreadyJoinSpace.Id, t)

	//join2
	wCh2 <- join(T_SPACE_ID, x, y, t)
	c2Rsp = getMsg(rCh2, t)
	unpack(pb.CmdJoinAck, c2Rsp, t)
	//user list
	c2Rsp = getMsg(rCh2, t)
	ulNtf = unpack(pb.CmdUserListNtf, c2Rsp, t).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID, t)
	//join_ntf
	c1Rsp = getMsg(rCh1, t)
	joinNtf := unpack(pb.CmdJoinNtf, c1Rsp, t).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID2, t)
}

func tLeave(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("LEAVE")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr, t)
	defer closeClient(wCh2)

	//未登陆时leave
	wCh1 <- leave(T_SPACE_ID, t)
	c1Rsp := getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)
	wCh1 <- leaveOpt(t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
	wCh2 <- login(T_APP_ID, T_USER_ID2, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)

	//未join时leave
	wCh1 <- leave(T_SPACE_ID, t) //不应该有msg包返回
	wCh1 <- leaveOpt(t)          //不应该有msg包返回

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//leave 未加入的space
	wCh1 <- leave("not_joined_space_id", t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//leave
	wCh1 <- leave(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLeaveAck, c1Rsp, t)
	//leave_ntf
	c2Rsp = getMsg(rCh2, t)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp, t).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)

	//leave opt
	wCh2 <- leaveOpt(t)
	c2Rsp = getMsg(rCh2, t)
	unpack(pb.CmdLeaveAck, c2Rsp, t)
}

func tMove(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("MOVE")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr, t)
	defer closeClient(wCh2)

	//未登录时move
	wCh1 <- move(0, 1, 1, t)
	c1Rsp := getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
	wCh2 <- login(T_APP_ID, T_USER_ID2, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)

	var x1, y1 float32 = 1, 2
	var x2, y2 float32 = (x1 + T_GRID_W), (y1 + T_GRID_H)
	var x3, y3 float32 = (x2 + T_GRID_W), (y2 + T_GRID_H)
	var time int32 = 0

	//未join时move
	wCh1 <- move(time, x1, y1, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//move
	time = 10
	wCh1 <- move(time, x1, y1, t)
	c1Rsp = getMsg(rCh1, t)
	moveAck := unpack(pb.CmdMoveAck, c1Rsp, t).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x1, t)
	assert(moveAck.GetPosY() == y1, t)
	//move_ntf
	c2Rsp = getMsg(rCh2, t)
	moveNtf := unpack(pb.CmdMoveNtf, c2Rsp, t).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x1, t)
	assert(moveNtf.GetPosY() == y1, t)

	//time ignore
	time = 9
	wCh1 <- move(time, x3, y3, t)

	//across move
	time = 11
	wCh1 <- move(time, x2, y2, t)
	c1Rsp = getMsg(rCh1, t)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp, t).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x2, t)
	assert(moveAck.GetPosY() == y2, t)
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//move_ntf
	c2Rsp = getMsg(rCh2, t)
	moveNtf = unpack(pb.CmdMoveNtf, c2Rsp, t).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x2, t)
	assert(moveNtf.GetPosY() == y2, t)

	//across move again
	time = 12
	wCh1 <- move(time, x3, y3, t)
	c1Rsp = getMsg(rCh1, t)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp, t).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x3, t)
	assert(moveAck.GetPosY() == y3, t)
	c1Rsp = getMsg(rCh1, t)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//move_ntf
	c2Rsp = getMsg(rCh2, t)
	moveNtf = unpack(pb.CmdMoveNtf, c2Rsp, t).(*pb.MoveNtf)
	assert(moveNtf.GetUserId() == T_USER_ID, t)
	assert(moveNtf.GetPosX() == x3, t)
	assert(moveNtf.GetPosY() == y3, t)
	//leave_ntf (离开了九宫格范围)
	c2Rsp = getMsg(rCh2, t)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp, t).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)

	//move back
	time = 13
	wCh1 <- move(time, x2, y2, t)
	c1Rsp = getMsg(rCh1, t)
	moveAck = unpack(pb.CmdMoveAck, c1Rsp, t).(*pb.MoveAck)
	assert(moveAck.GetPosX() == x2, t)
	assert(moveAck.GetPosY() == y2, t)
	c1Rsp = getMsg(rCh1, t)
	ulNtf = unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID2, t)
	//join_ntf (进入了九宫格范围)
	c2Rsp = getMsg(rCh2, t)
	joinNtf := unpack(pb.CmdJoinNtf, c2Rsp, t).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID, t)
	assert(joinNtf.GetUser().GetPosX() == x2, t)
	assert(joinNtf.GetUser().GetPosY() == y2, t)
}

func tBroadcast(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("BROADCAST")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)
	wCh2, rCh2 := startClient(tcpAddr, t)
	defer closeClient(wCh2)

	//未登录
	wCh1 <- broadcast(t)
	c1Rsp := getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
	wCh2 <- login(T_APP_ID, T_USER_ID2, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)

	//未join
	wCh1 <- broadcast(t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotJoinSpace.Id, t)

	//join
	joinSpace(rCh1, wCh1, rCh2, wCh2, t)

	//广播1
	wCh1 <- broadcast(t)
	c1Rsp = getMsg(rCh1, t)
	_ = unpack(pb.CmdBroadcastAck, c1Rsp, t).(*pb.BroadcastAck)
	//broadcast_ntf
	c2Rsp = getMsg(rCh2, t)
	bcNtf := unpack(pb.CmdBroadcastNtf, c2Rsp, t).(*pb.BroadcastNtf)
	assert(bcNtf.GetUserId() == T_USER_ID, t)

	//广播2
	//leave
	wCh1 <- leave(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLeaveAck, c1Rsp, t)
	//leave_ntf
	c2Rsp = getMsg(rCh2, t)
	leaveNtf := unpack(pb.CmdLeaveNtf, c2Rsp, t).(*pb.LeaveNtf)
	assert(leaveNtf.GetUserId() == T_USER_ID, t)
	//重新join
	//join
	wCh1 <- join(T_SPACE_ID, T_GRID_W*10, T_GRID_H*10, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//广播
	wCh1 <- broadcast(t)
	c1Rsp = getMsg(rCh1, t)
	_ = unpack(pb.CmdBroadcastAck, c1Rsp, t).(*pb.BroadcastAck)
	//none broadcast_ntf
	getNilMsg(rCh2, t)
}

func tQueryPos(tcpAddr string, t *testing.T) {
	fmt.Println("============================================")
	fmt.Println("QUERY_POS")
	fmt.Println("============================================")
	wCh1, rCh1 := startClient(tcpAddr, t)
	defer closeClient(wCh1)

	//自己未登录
	wCh1 <- broadcast(t)
	c1Rsp := getMsg(rCh1, t)
	eNtf := unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrNotLogin.Id, t)

	//login
	wCh1 <- login(T_APP_ID, T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)

	//不在线
	wCh1 <- querypos(T_USER_ID2, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrUserOffline.Id, t)

	//query
	wCh1 <- querypos(T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	qAck := unpack(pb.CmdQueryPosAck, c1Rsp, t).(*pb.QueryPosAck)
	assert(string(qAck.GetSpaceId()) == model.NIL, t)
	assert(qAck.GetPosX() == model.DEFAULT_POS_X, t)
	assert(qAck.GetPosY() == model.DEFAULT_POS_Y, t)
	assert(qAck.GetAngle() == model.DEFAULT_ANGLE, t)

	//join
	var x, y float32 = 190, 910
	wCh1 <- join(T_SPACE_ID, x, y, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//query
	wCh1 <- querypos(T_USER_ID, t)
	c1Rsp = getMsg(rCh1, t)
	qAck = unpack(pb.CmdQueryPosAck, c1Rsp, t).(*pb.QueryPosAck)
	assert(string(qAck.GetSpaceId()) == T_SPACE_ID, t)
	assert(qAck.GetPosX() == x, t)
	assert(qAck.GetPosY() == y, t)
	assert(qAck.GetAngle() == model.DEFAULT_ANGLE, t)
}

//两个user join
func joinSpace(rCh1, wCh1, rCh2, wCh2 chan []byte, t *testing.T) {
	//join
	wCh1 <- join(T_SPACE_ID, 0, 0, t)
	c1Rsp := getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//join2
	wCh2 <- join(T_SPACE_ID, 0, 0, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdJoinAck, c2Rsp, t)
	//user list
	c2Rsp = getMsg(rCh2, t)
	ulNtf = unpack(pb.CmdUserListNtf, c2Rsp, t).(*pb.UserListNtf)
	assert(1 == len(ulNtf.UserList), t)
	assert(ulNtf.UserList[0].GetUserId() == T_USER_ID, t)
	//join_ntf
	c1Rsp = getMsg(rCh1, t)
	joinNtf := unpack(pb.CmdJoinNtf, c1Rsp, t).(*pb.JoinNtf)
	assert(joinNtf.GetUser().GetUserId() == T_USER_ID2, t)
}

func startClient(tcpAddr string, t *testing.T) (chan []byte, chan []byte) {
	conn, err := net.Dial("tcp", tcpAddr)
	check(err, t)
	wCh := make(chan []byte)
	go writer(wCh, conn, t)
	rCh := make(chan []byte)
	go reader(rCh, conn, t)
	return wCh, rCh
}

func closeClient(wCh chan []byte) {
	close(wCh)
}

func getMsg(ch chan []byte, t *testing.T) []byte {
	select {
	case m := <-ch:
		return m
	case <-time.After(time.Second * 5):
		debug.PrintStack()
		t.Fatal("get channel msg timeout")
		return nil
	}
}
func getNilMsg(ch chan []byte, t *testing.T) {
	select {
	case m := <-ch:
		debug.PrintStack()
		t.Fatalf("get channel unexpect msg:%v", m)
	case <-time.After(time.Second * 5):
	}
}

func writer(ch chan []byte, conn net.Conn, t *testing.T) {
	const lenHead = 4
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		size := len(data)
		buf := make([]byte, size+lenHead)
		binary.BigEndian.PutUint32(buf[:lenHead], uint32(size))
		copy(buf[lenHead:], data)
		_, err := conn.Write(buf[:size+lenHead])
		check(err, t)
	}
	conn.Close()
}

func reader(ch chan []byte, conn net.Conn, t *testing.T) {
	buf := make([]byte, 256)
	const lenHead = 4
	for {
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		_, err := io.ReadAtLeast(conn, buf[:lenHead], lenHead)
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, t)
		size := int(binary.BigEndian.Uint32(buf[:lenHead]))
		// read data
		_, err = io.ReadAtLeast(conn, buf[:size], size)
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, t)
		ch <- buf[:size]
	}
}
