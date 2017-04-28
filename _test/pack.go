package _test

import (
	"github.com/golang/protobuf/proto"
	"jhqc.com/songcf/scene/pb"
	"runtime/debug"
	"testing"
)

func heartbeat(t *testing.T) []byte {
	p := &pb.HeartbeatReq{}
	return pack(pb.CmdHeartbeatReq, p, t)
}

func login(appId string, uid int32, t *testing.T) []byte {
	p := &pb.LoginReq{
		AppId:     []byte(appId),
		UserId:    &uid,
		UserToken: []byte("token"),
	}
	return pack(pb.CmdLoginReq, p, t)
}

func join(spaceId string, t *testing.T) []byte {
	useLast := true
	p := &pb.JoinReq{
		SpaceId: []byte(spaceId),
		UseLast: &useLast,
		ExData:  []byte(""),
	}
	return pack(pb.CmdJoinReq, p, t)
}

func leave(spaceId string, t *testing.T) []byte {
	p := &pb.LeaveReq{SpaceId: []byte(spaceId)}
	return pack(pb.CmdLeaveReq, p, t)
}
func leaveOpt(t *testing.T) []byte {
	p := &pb.LeaveReq{}
	return pack(pb.CmdLeaveReq, p, t)
}

func move(time int32, x, y float32, t *testing.T) []byte {
	var a float32 = 0
	p := &pb.MoveReq{
		PosX:  &x,
		PosY:  &y,
		Angle: &a,
		Time:  &time,
	}
	return pack(pb.CmdMoveReq, p, t)
}

func broadcast(t *testing.T) []byte {
	p := &pb.BroadcastReq{Data: []byte("test-broadcast-data")}
	return pack(pb.CmdBroadcastReq, p, t)
}

func querypos(uid int32, t *testing.T) []byte {
	p := &pb.QueryPosReq{UserId: &uid}
	return pack(pb.CmdQueryPosReq, p, t)
}

func pack(cmd int32, p proto.Message, t *testing.T) []byte {
	var vsn int32 = 1
	pl, err := proto.Marshal(p)
	check(err, t)
	req := &pb.Packet{
		Cmd:     &cmd,
		Vsn:     &vsn,
		Payload: pl,
	}
	data, err := proto.Marshal(req)
	check(err, t)
	return data
}

func unpack(cmd int32, data []byte, t *testing.T) proto.Message {
	var err error
	var p proto.Message
	// parse
	packet := &pb.Packet{}
	err = proto.Unmarshal(data, packet)
	check(err, t)
	cmd2 := packet.GetCmd()
	payload := packet.GetPayload()
	if cmd2 != cmd {
		debug.PrintStack()
		if cmd2 == pb.CmdErrorNtf {
			p = &pb.ErrorNtf{}
			err = proto.Unmarshal(payload, p)
			if err != nil {
				t.Fatalf("parse error_ntf failed, err=%v", err)
			}
			t.Fatalf("check pb error, error_ntf=%v\n", p)
		} else {
			t.Fatalf("check pb error, cmd=%v, packCmd=%v\n", cmd, cmd2)
		}
	}
	switch cmd {
	case pb.CmdHeartbeatAck:
		p = &pb.HeartbeatAck{}
	case pb.CmdLoginAck:
		p = &pb.LoginAck{}
	case pb.CmdJoinAck:
		p = &pb.JoinAck{}
	case pb.CmdJoinNtf:
		p = &pb.JoinNtf{}
	case pb.CmdLeaveAck:
		p = &pb.LeaveAck{}
	case pb.CmdLeaveNtf:
		p = &pb.LeaveNtf{}
	case pb.CmdMoveAck:
		p = &pb.MoveAck{}
	case pb.CmdMoveNtf:
		p = &pb.MoveNtf{}
	case pb.CmdBroadcastAck:
		p = &pb.BroadcastAck{}
	case pb.CmdBroadcastNtf:
		p = &pb.BroadcastNtf{}
	case pb.CmdQueryPosAck:
		p = &pb.QueryPosAck{}
	case pb.CmdErrorNtf:
		p = &pb.ErrorNtf{}
	case pb.CmdOfflineNtf:
		p = &pb.OfflineNtf{}
	case pb.CmdUserListNtf:
		p = &pb.UserListNtf{}
	default:
		t.Fatalf("check unknown cmd:%v", cmd)
	}
	err = proto.Unmarshal(payload, p)
	check(err, t)
	return p
}
