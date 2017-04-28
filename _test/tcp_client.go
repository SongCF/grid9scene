package _test

import (
	"encoding/binary"
	"fmt"
	"io"
	"jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"net"
	"strings"
	"testing"
	"time"
)

func testAllMsg(httpAddr, tcpAddr string, t *testing.T) {
	initAppSpace(httpAddr, t)

	tHeartbeat(tcpAddr, t)
	time.Sleep(time.Second)

	tLogin(tcpAddr, t)
	time.Sleep(time.Second)

	tJoin(tcpAddr, t)
	time.Sleep(time.Second)

	tLeave(tcpAddr, t)
	time.Sleep(time.Second)

	tMove(tcpAddr, t)
	time.Sleep(time.Second)

	tBroadcast(tcpAddr, t)
	time.Sleep(time.Second)

	tQueryPos(tcpAddr, t)
	time.Sleep(time.Second)
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
	assert(eNtf.GetCode() == pb.ErrDuplicateLogin.Id,
		fmt.Sprintf("duplicate login error ntf not match,e:%v", eNtf), t)
	//顶号登录
	wCh2 <- login(T_APP_ID, T_USER_ID, t)
	c2Rsp := getMsg(rCh2, t)
	unpack(pb.CmdLoginAck, c2Rsp, t)
	c1Rsp = getMsg(rCh1, t)
	offlineNtf := unpack(pb.CmdOfflineNtf, c1Rsp, t).(*pb.OfflineNtf)
	assert(offlineNtf.GetType() == model.OFFLINE_TYPE_OTHER_LOGIN,
		fmt.Sprintf("offline ntf error, %v", offlineNtf), t)
	//应用不存在
	wCh1 <- login("not_exist_app_id", T_USER_ID2, t)
	c1Rsp = getMsg(rCh1, t)
	eNtf = unpack(pb.CmdErrorNtf, c1Rsp, t).(*pb.ErrorNtf)
	assert(eNtf.GetCode() == pb.ErrAppNotExist.Id,
		fmt.Sprintf("app not exist login error ntf not match,e:%v", eNtf), t)
	//保持c1c2正常登陆
	wCh1 <- login(T_APP_ID, T_USER_ID2, t)
	c1Rsp = getMsg(rCh1, t)
	unpack(pb.CmdLoginAck, c1Rsp, t)
}

func tJoin(tcpAddr string, t *testing.T) {

}

func tLeave(tcpAddr string, t *testing.T) {

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
		t.Fatal("get channel msg timeout")
		return nil
	}
}

func writer(ch chan []byte, conn net.Conn, t *testing.T) {
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		size := len(data)
		buf := make([]byte, size+4)
		binary.BigEndian.PutUint32(buf, uint32(size))
		copy(buf[4:], data)
		_, err := conn.Write(buf[:size+4])
		check(err, t)
	}
	conn.Close()
}

func reader(ch chan []byte, conn net.Conn, t *testing.T) {
	buf := make([]byte, 256)
	for {
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		_, err := io.ReadAtLeast(conn, buf[:4], 4)
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, t)
		size := binary.BigEndian.Uint32(buf[:4])
		// read data
		_, err = io.ReadAtLeast(conn, buf[:size], int(size))
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, t)
		ch <- buf[:size]
	}
}
