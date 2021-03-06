package pb

import (
	"github.com/golang/protobuf/proto"
)

// 场景服务	10070000	10079999

type ErrInfo struct {
	Id   int32  `json:"code"`
	Desc string `json:"msg"`
	Ex   string `json:"response,omitempty"`
}

var (
	//无错
	ErrSuccess = &ErrInfo{Id: 10000000, Desc: "Success"}

	//数据库出错
	ErrQueryDBError = &ErrInfo{Id: 10070001, Desc: "query database error"}
	//请求超时
	ErrTimeOut = &ErrInfo{Id: 10070002, Desc: "request timeout"}
	//服务器忙
	ErrServerBusy = &ErrInfo{Id: 10070003, Desc: "Server Busy"}
	//未登陆
	ErrNotLogin = &ErrInfo{Id: 10070101, Desc: "has not login"}
	//重复登陆
	ErrDuplicateLogin = &ErrInfo{Id: 10070102, Desc: "duplicate login"}
	//%% 指定的玩家id未找到
	//user_not_found()-> {10070103, <<"user not found">>}.
	//错误的user
	ErrUser = &ErrInfo{Id: 10070104, Desc: "error user"}
	//用户不在线
	ErrUserOffline = &ErrInfo{Id: 10070105, Desc: "User offline"}
	//
	//
	//应用不存在
	ErrAppNotExist = &ErrInfo{Id: 10070201, Desc: "app not exist"}
	//应用已存在
	ErrAppAlreadyExist = &ErrInfo{Id: 10070202, Desc: "app already exist"}
	//场景不存在
	ErrSpaceNotExist = &ErrInfo{Id: 10070301, Desc: "space not exist"}
	//场景已经存在
	ErrSpaceAlreadyExist = &ErrInfo{Id: 10070302, Desc: "space already exist"}
	//缺少space
	ErrMissSpaceId = &ErrInfo{Id: 10070303, Desc: "miss space id"}
	//已加入该场景了
	ErrAlreadyJoinSpace = &ErrInfo{Id: 10070304, Desc: "already join current space"}
	//未加入任何场景
	ErrNotJoinSpace = &ErrInfo{Id: 10070305, Desc: "user not join some space"}
	//
	//
	//消息格式错误
	ErrMsgFormat = &ErrInfo{Id: 10070401, Desc: "error message format"}
	//%% 未知消息
	//unknown_msg()-> {10070402, <<"unknown message">>}.

	//不支持的协议/方法
	ErrCmdNotSupport = &ErrInfo{Id: 10070404, Desc: "cmd not support"}
)

// packaging error
func Error(req int32, err *ErrInfo) (int32, proto.Message) {
	e := &ErrorNtf{
		Code: &err.Id,
		Msg:  []byte(err.Desc),
		Req:  &req,
	}
	return CmdErrorNtf, e
}
