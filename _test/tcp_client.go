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

	//未登陆时join
	wCh1 <- join(T_SPACE_ID, t)
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
	wCh1 <- join("not_exist_space_id", t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrSpaceNotExist.Id, t)

	//join
	wCh1 <- join(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)

	//重复join
	wCh1 <- join(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAlreadyJoinSpace.Id, t)

	//join2
	wCh2 <- join(T_SPACE_ID, t)
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
	wCh1 <- join(T_SPACE_ID, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdJoinAck, c1Rsp, t)
	//user list
	c1Rsp = getMsg(rCh1, t)
	ulNtf := unpack(pb.CmdUserListNtf, c1Rsp, t).(*pb.UserListNtf)
	assert(0 == len(ulNtf.UserList), t)
	//join2
	wCh2 <- join(T_SPACE_ID, t)
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

}

func tBroadcast(tcpAddr string, t *testing.T) {

}

func tQueryPos(tcpAddr string, t *testing.T) {

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
	case <-time.After(time.Second * 2):
		debug.PrintStack()
		t.Fatal("get channel msg timeout")
		return nil
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
