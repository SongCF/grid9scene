package _test

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/SongCF/scene/pb"
	"runtime/debug"
)

func heartbeat() []byte {
	p := &pb.HeartbeatReq{}
	return pack(pb.CmdHeartbeatReq, p)
}

func login(appId string, uid int32) []byte {
	p := &pb.LoginReq{
		AppId:     []byte(appId),
		UserId:    &uid,
		UserToken: []byte("token"),
	}
	return pack(pb.CmdLoginReq, p)
}

func join(spaceId string, x, y float32) []byte {
	useLast := false
	p := &pb.JoinReq{
		SpaceId: []byte(spaceId),
		PosX:    &x,
		PosY:    &y,
		UseLast: &useLast,
		ExData:  []byte(""),
	}
	return pack(pb.CmdJoinReq, p)
}
func joinLastPos(spaceId string) []byte {
	useLast := true
	p := &pb.JoinReq{
		SpaceId: []byte(spaceId),
		UseLast: &useLast,
		ExData:  []byte(""),
	}
	return pack(pb.CmdJoinReq, p)
}

func leave(spaceId string) []byte {
	p := &pb.LeaveReq{SpaceId: []byte(spaceId)}
	return pack(pb.CmdLeaveReq, p)
}
func leaveOpt() []byte {
	p := &pb.LeaveReq{}
	return pack(pb.CmdLeaveReq, p)
}

func move(time int32, x, y float32) []byte {
	var a float32 = 0
	p := &pb.MoveReq{
		PosX:  &x,
		PosY:  &y,
		Angle: &a,
		Time:  &time,
	}
	return pack(pb.CmdMoveReq, p)
}

func broadcast() []byte {
	p := &pb.BroadcastReq{Data: []byte("test-broadcast-data")}
	return pack(pb.CmdBroadcastReq, p)
}

func querypos(uid int32) []byte {
	p := &pb.QueryPosReq{UserId: &uid}
	return pack(pb.CmdQueryPosReq, p)
}

func pack(cmd int32, p proto.Message) []byte {
	var vsn int32 = 1
	pl, err := proto.Marshal(p)
	check(err, "pack marshal payload failed:")
	req := &pb.Packet{
		Cmd:     &cmd,
		Vsn:     &vsn,
		Payload: pl,
	}
	data, err := proto.Marshal(req)
	check(err, "pack marshal packet failed:")
	return data
}

func unpack(cmd int32, data []byte) proto.Message {
	var err error
	var p proto.Message
	// parse
	packet := &pb.Packet{}
	err = proto.Unmarshal(data, packet)
	check(err, "unpack unmarshal packet failed:")
	cmd2 := packet.GetCmd()
	payload := packet.GetPayload()
	if cmd2 != cmd {
		debug.PrintStack()
		if cmd2 == pb.CmdErrorNtf {
			p = &pb.ErrorNtf{}
			err = proto.Unmarshal(payload, p)
			check(err, "parse error_ntf failed")
			err = errors.New(fmt.Sprintf("check pb error, error_ntf=%v\n", p))
			check(err, "")
		} else {
			err = errors.New(fmt.Sprintf("check pb error, cmd=%v, packCmd=%v\n", cmd, cmd2))
			check(err, "")
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
		err = errors.New(fmt.Sprintf("check unknown cmd:%v", cmd))
		check(err, "")
	}
	err = proto.Unmarshal(payload, p)
	check(err, "unpack unmarshal payload failed:")
	return p
}
