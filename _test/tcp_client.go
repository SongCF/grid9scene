package _test

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"io"
	"jhqc.com/songcf/scene/pb"
	"net"
	"testing"
	"time"
)

var vsn int32 = 1

func testAllMsg(httpAddr, tcpAddr string, t *testing.T) {
	initAppSpace(httpAddr, t)

	// start client 1
	conn1, err := net.Dial("tcp", tcpAddr)
	check(err, t)
	c1WCh := make(chan []byte)
	go writer(c1WCh, conn1, t)
	c1RCh := make(chan *pb.Packet)
	go reader(c1RCh, conn1, t)

	//// start client 2
	//conn2, err := net.Dial("tcp", tcpAddr)
	//check(err, t)
	//c2WCh := make(chan []byte)
	//go writer(c2WCh, conn2, t)
	//c2RCh := make(chan *pb.Packet)
	//go reader(c2RCh, conn2, t)

	//
	c1WCh <- login(T_USER_ID, t)
	pack := <-c1RCh
	checkPb(pb.CmdLoginAck, pack, t)
}

func writer(ch chan []byte, conn net.Conn, t *testing.T) {
	for {
		data := <-ch
		size := len(data)
		buf := make([]byte, size+4)
		binary.BigEndian.PutUint32(buf, uint32(size))
		copy(buf[4:], data)
		_, err := conn.Write(buf[:size+4])
		check(err, t)
	}
}

func reader(ch chan *pb.Packet, conn net.Conn, t *testing.T) {
	buf := make([]byte, 256)
	for {
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		_, err := io.ReadAtLeast(conn, buf[:4], 4)
		check(err, t)
		size := binary.BigEndian.Uint32(buf[:4])
		// read data
		_, err = io.ReadAtLeast(conn, buf[:size], int(size))
		check(err, t)
		// parse
		packet := &pb.Packet{}
		err = proto.Unmarshal(buf[:size], packet)
		check(err, t)
		ch <- packet
	}
}

func heartbeat() {

}

func login(uid int32, t *testing.T) []byte {
	p := &pb.LoginReq{
		AppId:     []byte(T_APP_ID),
		UserId:    &uid,
		UserToken: []byte("token"),
	}
	pl, err := proto.Marshal(p)
	check(err, t)
	var cmd int32 = pb.CmdLoginReq
	req := &pb.Packet{
		Cmd:     &cmd,
		Vsn:     &vsn,
		Payload: pl,
	}
	data, err := proto.Marshal(req)
	check(err, t)
	return data
}

func checkPb(cmd int32, pack *pb.Packet, t *testing.T) {
	var err error
	var p proto.Message
	data := pack.GetPayload()
	if pack.GetCmd() == pb.CmdErrorNtf {
		p = &pb.ErrorNtf{}
		err = proto.Unmarshal(data, p)
		if err != nil {
			t.Fatalf("parse error_ntf failed, err=%v", err)
		}
		t.Fatalf("checkPB error, error_ntf=%v\n", p)
	}
	if pack.GetCmd() != cmd {
		t.Fatalf("checkPB error, cmd=%v, packCmd=%v\n", cmd, pack.GetCmd())
	}
	switch cmd {
	case pb.CmdHeartbeatAck:
		p = &pb.HeartbeatAck{}
	case pb.CmdLoginAck:
		p = &pb.LoginAck{}
	case pb.CmdJoinAck:
	case pb.CmdJoinNtf:
	case pb.CmdLeaveAck:
	case pb.CmdLeaveNtf:
	case pb.CmdMoveAck:
	case pb.CmdMoveNtf:
	case pb.CmdBroadcastAck:
	case pb.CmdBroadcastNtf:
	case pb.CmdQueryPosAck:
	default:
		t.Fatalf("check unknown cmd:%v", cmd)
	}
	err = proto.Unmarshal(data, p)
	check(err, t)
}
