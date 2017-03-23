package pb


// 场景服务	10070000	10079999

type ErrInfo struct {
	Id int
	Desc string
}

const (
	ErrTimeOut = &ErrInfo{Id:10070002, Desc:"request timeout"} //服务器错误
	ErrNotLogin = &ErrInfo{Id:10070101, Desc:"has not login"} //未登录
	//%% 重复登陆
	//duplicate_login()-> {10070102, <<"duplicate login">>}.
	//%% 指定的玩家id未找到
	//user_not_found()-> {10070103, <<"user not found">>}.
	//%% 错误的user
	//error_user()-> {10070104, <<"error user">>}.
	//
	//
	//%% 应用不存在
	//app_not_exist()-> {10070201, <<"app not exist">>}.
	//
	//
	//%% 场景不存在
	//space_not_exist()-> {10070301, <<"space not exist">>}.
	//%% 场景已经存在
	//space_already_exist()-> {10070302, <<"space already exist">>}.
	//%% 格子宽高不正确
	//space_grid_error()-> {10070303, <<"space grid size error">>}.
	//
	//
	//消息格式错误
	ErrMsgFormat = &ErrInfo{Id:10070401, Desc:"error message format"}
	//%% 未知消息
	//unknown_msg()-> {10070402, <<"unknown message">>}.
	//%% 缺少参数 Param::binary()
	//miss_param(Param) when is_atom(Param) -> miss_param(atom_to_binary(Param,utf8));
	//miss_param(Param)-> {10070403, <<"Missing parameter: ", Param/binary>>}.

	//不支持的协议/方法
	ErrCmdNotSupport = &ErrInfo{Id:10070404, Desc:"cmd not support"}

)
