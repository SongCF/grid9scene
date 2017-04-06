// Code generated by protoc-gen-go.
// source: payload.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	payload.proto

It has these top-level messages:
	Packet
	UserData
	ErrorNtf
	LoginReq
	LoginAck
	OfflineNtf
	JoinReq
	JoinAck
	JoinNtf
	UserListNtf
	LeaveReq
	LeaveAck
	LeaveNtf
	MoveReq
	MoveAck
	MoveNtf
	BroadcastReq
	BroadcastAck
	BroadcastNtf
	QueryPosReq
	QueryPosAck
	HeartbeatReq
	HeartbeatAck
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// 第一层pb
type Packet struct {
	Cmd              *int32 `protobuf:"varint,1,req,name=cmd" json:"cmd,omitempty"`
	Vsn              *int32 `protobuf:"varint,2,req,name=vsn" json:"vsn,omitempty"`
	Payload          []byte `protobuf:"bytes,3,req,name=payload" json:"payload,omitempty"`
	Service          []byte `protobuf:"bytes,4,opt,name=service" json:"service,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Packet) Reset()                    { *m = Packet{} }
func (m *Packet) String() string            { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()               {}
func (*Packet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Packet) GetCmd() int32 {
	if m != nil && m.Cmd != nil {
		return *m.Cmd
	}
	return 0
}

func (m *Packet) GetVsn() int32 {
	if m != nil && m.Vsn != nil {
		return *m.Vsn
	}
	return 0
}

func (m *Packet) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Packet) GetService() []byte {
	if m != nil {
		return m.Service
	}
	return nil
}

// 玩家
type UserData struct {
	UserId           *int32   `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	PosX             *float32 `protobuf:"fixed32,2,req,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,3,req,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,4,req,name=angle" json:"angle,omitempty"`
	ExData           []byte   `protobuf:"bytes,5,opt,name=ex_data,json=exData" json:"ex_data,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *UserData) Reset()                    { *m = UserData{} }
func (m *UserData) String() string            { return proto.CompactTextString(m) }
func (*UserData) ProtoMessage()               {}
func (*UserData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserData) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *UserData) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *UserData) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *UserData) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

func (m *UserData) GetExData() []byte {
	if m != nil {
		return m.ExData
	}
	return nil
}

// 通用错误
// cmd:10012
type ErrorNtf struct {
	Code             *int32 `protobuf:"varint,1,req,name=code" json:"code,omitempty"`
	Msg              []byte `protobuf:"bytes,2,req,name=msg" json:"msg,omitempty"`
	Req              *int32 `protobuf:"varint,3,req,name=req" json:"req,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ErrorNtf) Reset()                    { *m = ErrorNtf{} }
func (m *ErrorNtf) String() string            { return proto.CompactTextString(m) }
func (*ErrorNtf) ProtoMessage()               {}
func (*ErrorNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ErrorNtf) GetCode() int32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *ErrorNtf) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *ErrorNtf) GetReq() int32 {
	if m != nil && m.Req != nil {
		return *m.Req
	}
	return 0
}

// 用于登录服务器 验证
// cmd:10020
type LoginReq struct {
	AppId            []byte `protobuf:"bytes,1,req,name=app_id,json=appId" json:"app_id,omitempty"`
	UserId           *int32 `protobuf:"varint,2,req,name=user_id,json=userId" json:"user_id,omitempty"`
	UserToken        []byte `protobuf:"bytes,4,req,name=user_token,json=userToken" json:"user_token,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *LoginReq) Reset()                    { *m = LoginReq{} }
func (m *LoginReq) String() string            { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()               {}
func (*LoginReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LoginReq) GetAppId() []byte {
	if m != nil {
		return m.AppId
	}
	return nil
}

func (m *LoginReq) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *LoginReq) GetUserToken() []byte {
	if m != nil {
		return m.UserToken
	}
	return nil
}

// 用于登录服务器 验证 返回
// cmd:10021
type LoginAck struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *LoginAck) Reset()                    { *m = LoginAck{} }
func (m *LoginAck) String() string            { return proto.CompactTextString(m) }
func (*LoginAck) ProtoMessage()               {}
func (*LoginAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// 推送下线消息
// cmd:10032
type OfflineNtf struct {
	Type             *int32 `protobuf:"varint,1,req,name=type" json:"type,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *OfflineNtf) Reset()                    { *m = OfflineNtf{} }
func (m *OfflineNtf) String() string            { return proto.CompactTextString(m) }
func (*OfflineNtf) ProtoMessage()               {}
func (*OfflineNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *OfflineNtf) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

// 加入场景请求
// cmd:20010
type JoinReq struct {
	SpaceId          []byte   `protobuf:"bytes,1,req,name=space_id,json=spaceId" json:"space_id,omitempty"`
	PosX             *float32 `protobuf:"fixed32,2,opt,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,3,opt,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,4,opt,name=angle" json:"angle,omitempty"`
	UseLast          *bool    `protobuf:"varint,5,req,name=use_last,json=useLast" json:"use_last,omitempty"`
	ExData           []byte   `protobuf:"bytes,6,req,name=ex_data,json=exData" json:"ex_data,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *JoinReq) Reset()                    { *m = JoinReq{} }
func (m *JoinReq) String() string            { return proto.CompactTextString(m) }
func (*JoinReq) ProtoMessage()               {}
func (*JoinReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *JoinReq) GetSpaceId() []byte {
	if m != nil {
		return m.SpaceId
	}
	return nil
}

func (m *JoinReq) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *JoinReq) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *JoinReq) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

func (m *JoinReq) GetUseLast() bool {
	if m != nil && m.UseLast != nil {
		return *m.UseLast
	}
	return false
}

func (m *JoinReq) GetExData() []byte {
	if m != nil {
		return m.ExData
	}
	return nil
}

// 加入场景返回
// cmd:20011
type JoinAck struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *JoinAck) Reset()                    { *m = JoinAck{} }
func (m *JoinAck) String() string            { return proto.CompactTextString(m) }
func (*JoinAck) ProtoMessage()               {}
func (*JoinAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

// 加入场景推送
// cmd:20012
type JoinNtf struct {
	User             *UserData `protobuf:"bytes,1,req,name=user" json:"user,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *JoinNtf) Reset()                    { *m = JoinNtf{} }
func (m *JoinNtf) String() string            { return proto.CompactTextString(m) }
func (*JoinNtf) ProtoMessage()               {}
func (*JoinNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *JoinNtf) GetUser() *UserData {
	if m != nil {
		return m.User
	}
	return nil
}

// 推送其它格子的玩家列表
// cmd:20013
type UserListNtf struct {
	UserList         []*UserData `protobuf:"bytes,1,rep,name=user_list,json=userList" json:"user_list,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *UserListNtf) Reset()                    { *m = UserListNtf{} }
func (m *UserListNtf) String() string            { return proto.CompactTextString(m) }
func (*UserListNtf) ProtoMessage()               {}
func (*UserListNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *UserListNtf) GetUserList() []*UserData {
	if m != nil {
		return m.UserList
	}
	return nil
}

// 离开场景请求
// cmd:20020
type LeaveReq struct {
	SpaceId          []byte `protobuf:"bytes,1,opt,name=space_id,json=spaceId" json:"space_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *LeaveReq) Reset()                    { *m = LeaveReq{} }
func (m *LeaveReq) String() string            { return proto.CompactTextString(m) }
func (*LeaveReq) ProtoMessage()               {}
func (*LeaveReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *LeaveReq) GetSpaceId() []byte {
	if m != nil {
		return m.SpaceId
	}
	return nil
}

// 离开场景返回
// cmd:20021
type LeaveAck struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *LeaveAck) Reset()                    { *m = LeaveAck{} }
func (m *LeaveAck) String() string            { return proto.CompactTextString(m) }
func (*LeaveAck) ProtoMessage()               {}
func (*LeaveAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

// 离开场景推送
// cmd:20022
type LeaveNtf struct {
	UserId           *int32 `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *LeaveNtf) Reset()                    { *m = LeaveNtf{} }
func (m *LeaveNtf) String() string            { return proto.CompactTextString(m) }
func (*LeaveNtf) ProtoMessage()               {}
func (*LeaveNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *LeaveNtf) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 场景移动请求
// cmd:20030
type MoveReq struct {
	PosX             *float32 `protobuf:"fixed32,1,req,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,2,req,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,3,req,name=angle" json:"angle,omitempty"`
	Time             *int32   `protobuf:"varint,4,req,name=time" json:"time,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MoveReq) Reset()                    { *m = MoveReq{} }
func (m *MoveReq) String() string            { return proto.CompactTextString(m) }
func (*MoveReq) ProtoMessage()               {}
func (*MoveReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *MoveReq) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *MoveReq) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *MoveReq) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

func (m *MoveReq) GetTime() int32 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

// 场景移动返回
// cmd:20031
type MoveAck struct {
	PosX             *float32 `protobuf:"fixed32,1,req,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,2,req,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,3,req,name=angle" json:"angle,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MoveAck) Reset()                    { *m = MoveAck{} }
func (m *MoveAck) String() string            { return proto.CompactTextString(m) }
func (*MoveAck) ProtoMessage()               {}
func (*MoveAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *MoveAck) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *MoveAck) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *MoveAck) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

// 场景移动推送
// cmd:20032
type MoveNtf struct {
	UserId           *int32   `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	PosX             *float32 `protobuf:"fixed32,2,req,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,3,req,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,4,req,name=angle" json:"angle,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MoveNtf) Reset()                    { *m = MoveNtf{} }
func (m *MoveNtf) String() string            { return proto.CompactTextString(m) }
func (*MoveNtf) ProtoMessage()               {}
func (*MoveNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *MoveNtf) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *MoveNtf) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *MoveNtf) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *MoveNtf) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

// 场景广播请求
// cmd:20040
type BroadcastReq struct {
	Data             []byte `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *BroadcastReq) Reset()                    { *m = BroadcastReq{} }
func (m *BroadcastReq) String() string            { return proto.CompactTextString(m) }
func (*BroadcastReq) ProtoMessage()               {}
func (*BroadcastReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *BroadcastReq) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// 场景广播返回
// cmd:20041
type BroadcastAck struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *BroadcastAck) Reset()                    { *m = BroadcastAck{} }
func (m *BroadcastAck) String() string            { return proto.CompactTextString(m) }
func (*BroadcastAck) ProtoMessage()               {}
func (*BroadcastAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

// 场景广播推送
// cmd:20042
type BroadcastNtf struct {
	UserId           *int32 `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	Data             []byte `protobuf:"bytes,2,req,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *BroadcastNtf) Reset()                    { *m = BroadcastNtf{} }
func (m *BroadcastNtf) String() string            { return proto.CompactTextString(m) }
func (*BroadcastNtf) ProtoMessage()               {}
func (*BroadcastNtf) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *BroadcastNtf) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *BroadcastNtf) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// 查询位置请求
// cmd:20050
type QueryPosReq struct {
	UserId           *int32 `protobuf:"varint,1,req,name=user_id,json=userId" json:"user_id,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *QueryPosReq) Reset()                    { *m = QueryPosReq{} }
func (m *QueryPosReq) String() string            { return proto.CompactTextString(m) }
func (*QueryPosReq) ProtoMessage()               {}
func (*QueryPosReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *QueryPosReq) GetUserId() int32 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

// 查询位置返回
// cmd:20051
type QueryPosAck struct {
	SpaceId          []byte   `protobuf:"bytes,1,req,name=space_id,json=spaceId" json:"space_id,omitempty"`
	PosX             *float32 `protobuf:"fixed32,2,req,name=pos_x,json=posX" json:"pos_x,omitempty"`
	PosY             *float32 `protobuf:"fixed32,3,req,name=pos_y,json=posY" json:"pos_y,omitempty"`
	Angle            *float32 `protobuf:"fixed32,4,req,name=angle" json:"angle,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *QueryPosAck) Reset()                    { *m = QueryPosAck{} }
func (m *QueryPosAck) String() string            { return proto.CompactTextString(m) }
func (*QueryPosAck) ProtoMessage()               {}
func (*QueryPosAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *QueryPosAck) GetSpaceId() []byte {
	if m != nil {
		return m.SpaceId
	}
	return nil
}

func (m *QueryPosAck) GetPosX() float32 {
	if m != nil && m.PosX != nil {
		return *m.PosX
	}
	return 0
}

func (m *QueryPosAck) GetPosY() float32 {
	if m != nil && m.PosY != nil {
		return *m.PosY
	}
	return 0
}

func (m *QueryPosAck) GetAngle() float32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

// 系统预留，用于表示心跳
// cmd:20061
type HeartbeatReq struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *HeartbeatReq) Reset()                    { *m = HeartbeatReq{} }
func (m *HeartbeatReq) String() string            { return proto.CompactTextString(m) }
func (*HeartbeatReq) ProtoMessage()               {}
func (*HeartbeatReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

// cmd:20062
type HeartbeatAck struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *HeartbeatAck) Reset()                    { *m = HeartbeatAck{} }
func (m *HeartbeatAck) String() string            { return proto.CompactTextString(m) }
func (*HeartbeatAck) ProtoMessage()               {}
func (*HeartbeatAck) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

func init() {
	proto.RegisterType((*Packet)(nil), "pb.packet")
	proto.RegisterType((*UserData)(nil), "pb.user_data")
	proto.RegisterType((*ErrorNtf)(nil), "pb.error_ntf")
	proto.RegisterType((*LoginReq)(nil), "pb.login_req")
	proto.RegisterType((*LoginAck)(nil), "pb.login_ack")
	proto.RegisterType((*OfflineNtf)(nil), "pb.offline_ntf")
	proto.RegisterType((*JoinReq)(nil), "pb.join_req")
	proto.RegisterType((*JoinAck)(nil), "pb.join_ack")
	proto.RegisterType((*JoinNtf)(nil), "pb.join_ntf")
	proto.RegisterType((*UserListNtf)(nil), "pb.user_list_ntf")
	proto.RegisterType((*LeaveReq)(nil), "pb.leave_req")
	proto.RegisterType((*LeaveAck)(nil), "pb.leave_ack")
	proto.RegisterType((*LeaveNtf)(nil), "pb.leave_ntf")
	proto.RegisterType((*MoveReq)(nil), "pb.move_req")
	proto.RegisterType((*MoveAck)(nil), "pb.move_ack")
	proto.RegisterType((*MoveNtf)(nil), "pb.move_ntf")
	proto.RegisterType((*BroadcastReq)(nil), "pb.broadcast_req")
	proto.RegisterType((*BroadcastAck)(nil), "pb.broadcast_ack")
	proto.RegisterType((*BroadcastNtf)(nil), "pb.broadcast_ntf")
	proto.RegisterType((*QueryPosReq)(nil), "pb.query_pos_req")
	proto.RegisterType((*QueryPosAck)(nil), "pb.query_pos_ack")
	proto.RegisterType((*HeartbeatReq)(nil), "pb.heartbeat_req")
	proto.RegisterType((*HeartbeatAck)(nil), "pb.heartbeat_ack")
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 558 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x95, 0x9d, 0x38, 0x4d, 0xa6, 0xb1, 0x8a, 0x0c, 0x08, 0xf7, 0x80, 0x94, 0x2e, 0x08, 0x45,
	0x48, 0xf4, 0xc0, 0x15, 0x6e, 0x70, 0x29, 0xea, 0xc9, 0xe2, 0x00, 0x1c, 0x30, 0x13, 0x7b, 0x13,
	0xdc, 0x38, 0xde, 0xcd, 0xee, 0x26, 0x4a, 0x0e, 0x7c, 0x08, 0x7f, 0x8b, 0x66, 0x62, 0xa7, 0x8e,
	0x48, 0x2b, 0x44, 0x6f, 0x33, 0x6f, 0xd7, 0xf3, 0xde, 0xbc, 0x7d, 0x86, 0x50, 0xe3, 0xb6, 0x54,
	0x98, 0x5f, 0x6a, 0xa3, 0x9c, 0x8a, 0x7c, 0x3d, 0x11, 0x3f, 0xa0, 0xa7, 0x31, 0x9b, 0x4b, 0x17,
	0x3d, 0x82, 0x4e, 0xb6, 0xc8, 0x63, 0x6f, 0xe4, 0x8f, 0x83, 0x84, 0x4a, 0x42, 0xd6, 0xb6, 0x8a,
	0xfd, 0x1d, 0xb2, 0xb6, 0x55, 0x14, 0xc3, 0x49, 0x3d, 0x22, 0xee, 0x8c, 0xfc, 0xf1, 0x30, 0x69,
	0x5a, 0x3a, 0xb1, 0xd2, 0xac, 0x8b, 0x4c, 0xc6, 0xdd, 0x91, 0x47, 0x27, 0x75, 0x2b, 0x7e, 0xc1,
	0x60, 0x65, 0xa5, 0x49, 0x73, 0x74, 0x18, 0x3d, 0x83, 0x13, 0x6e, 0x8a, 0x86, 0xa8, 0x47, 0xed,
	0x55, 0x1e, 0x3d, 0x86, 0x40, 0x2b, 0x9b, 0x6e, 0x98, 0xcd, 0x4f, 0xba, 0x5a, 0xd9, 0x2f, 0x0d,
	0xb8, 0x65, 0xb2, 0x1d, 0xf8, 0x35, 0x7a, 0x02, 0x01, 0x56, 0xb3, 0x92, 0x78, 0x08, 0xdc, 0x35,
	0x34, 0x58, 0x6e, 0x98, 0x23, 0x0e, 0x98, 0xbf, 0x27, 0x37, 0x1f, 0xd1, 0xa1, 0xf8, 0x00, 0x03,
	0x69, 0x8c, 0x32, 0x69, 0xe5, 0xa6, 0x51, 0x04, 0xdd, 0x4c, 0xe5, 0xb2, 0xe6, 0xe6, 0x9a, 0xb6,
	0x5c, 0xd8, 0x19, 0xf3, 0x0e, 0x13, 0x2a, 0x09, 0x31, 0x72, 0xc9, 0xa4, 0x41, 0x42, 0xa5, 0xf8,
	0x06, 0x83, 0x52, 0xcd, 0x8a, 0x2a, 0x35, 0x72, 0x19, 0x3d, 0x85, 0x1e, 0x6a, 0xdd, 0xac, 0x30,
	0x4c, 0x02, 0xd4, 0xfa, 0x2a, 0x6f, 0xaf, 0xe6, 0x1f, 0xac, 0xf6, 0x1c, 0x80, 0x0f, 0x9c, 0x9a,
	0xcb, 0x8a, 0x55, 0x0f, 0x13, 0xb6, 0xe4, 0x33, 0x01, 0xe2, 0xb4, 0x99, 0x8d, 0xd9, 0x5c, 0x5c,
	0xc0, 0xa9, 0x9a, 0x4e, 0xcb, 0xa2, 0x92, 0x8d, 0x5e, 0xb7, 0xd5, 0x7b, 0xbd, 0x54, 0x8b, 0xdf,
	0x1e, 0xf4, 0x6f, 0x54, 0xad, 0xe5, 0x1c, 0xfa, 0x56, 0x63, 0x26, 0x6f, 0xd5, 0x9c, 0x70, 0x7f,
	0xe8, 0xa8, 0x77, 0xcc, 0x51, 0xef, 0x98, 0xa3, 0xde, 0xad, 0xa3, 0xe7, 0xd0, 0x5f, 0x59, 0x99,
	0x96, 0x68, 0x5d, 0x1c, 0x8c, 0xfc, 0x71, 0x3f, 0xa1, 0xfd, 0xae, 0xd1, 0xba, 0xb6, 0xd9, 0x3d,
	0x26, 0x6d, 0xcc, 0x86, 0x5a, 0x1a, 0xad, 0xf2, 0xa6, 0xae, 0x69, 0x8f, 0x0b, 0xe8, 0xd2, 0xc2,
	0x2c, 0xf1, 0xf4, 0x6d, 0x78, 0xa9, 0x27, 0x97, 0xfb, 0x4c, 0x24, 0x7c, 0x24, 0xde, 0x41, 0xc8,
	0x50, 0x59, 0x58, 0xc7, 0xdf, 0xbc, 0xae, 0x73, 0x43, 0x40, 0xec, 0x8d, 0x3a, 0x7f, 0x7f, 0x48,
	0xfa, 0xcc, 0x75, 0x61, 0x9d, 0x78, 0x05, 0x83, 0x52, 0xe2, 0x5a, 0x1e, 0xf1, 0xc4, 0x6b, 0x79,
	0xc2, 0x5e, 0xf3, 0x3d, 0x12, 0xf8, 0xb2, 0x69, 0x88, 0xed, 0xae, 0x60, 0x8a, 0xef, 0xd0, 0x5f,
	0xa8, 0x7a, 0xf2, 0xde, 0x52, 0xef, 0x58, 0x48, 0xfd, 0x63, 0x21, 0xed, 0xb4, 0x43, 0x4a, 0xcf,
	0x59, 0x2c, 0x76, 0xc9, 0xa5, 0xe7, 0x2c, 0x16, 0x52, 0x7c, 0xaa, 0xe7, 0x63, 0x36, 0x7f, 0xe8,
	0x7c, 0x91, 0xd5, 0xb3, 0xee, 0x5b, 0xe8, 0xa1, 0x7f, 0x9a, 0x78, 0x01, 0xe1, 0xc4, 0x28, 0xcc,
	0x33, 0xb4, 0x8e, 0x5d, 0x89, 0xa0, 0xcb, 0x51, 0xd8, 0xe5, 0x8f, 0x6b, 0x71, 0xd6, 0xbe, 0x44,
	0x66, 0xbf, 0x6f, 0x03, 0xf7, 0xea, 0x6b, 0xc6, 0xf9, 0xad, 0x71, 0x63, 0x08, 0x97, 0x2b, 0x69,
	0xb6, 0x29, 0x89, 0x24, 0xce, 0x3b, 0x9f, 0xeb, 0xa6, 0x7d, 0x93, 0x3c, 0xfd, 0xb7, 0x3f, 0xe4,
	0xbf, 0x9c, 0x38, 0x83, 0xf0, 0xa7, 0x44, 0xe3, 0x26, 0x12, 0xd9, 0x89, 0x43, 0x00, 0xb3, 0xf9,
	0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0xab, 0x10, 0x02, 0x71, 0x05, 0x00, 0x00,
}
