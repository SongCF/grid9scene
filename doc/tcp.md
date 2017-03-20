场景服务TCP协议
==========

TCP以google protocol buffer作为协议格式。TCP数据包中，固定前4字节为后面数据包的长度

包头 (4bytes) | 数据段（长度由包头指定）
---|---

后面的数据包pb格式为：
```
message packet {
  required int32 cmd = 1;     //各接口的消息编号
  required int32 vsn = 2;     //接口版本号，默认填1
  required bytes payload = 3; //protobuf
  optional bytes service = 4;
}
```
其中`payload`字段为各个接口自己的协议pb。

---

# 用户登录
同一用户，如果在两个地方同时登陆，先登录的会被挤下线，并收到`offline_ntf`消息推送
```
// 用于登录服务器 验证
// cmd:10020
message login_req {
    required bytes app_id = 1;
    required int32 user_id = 2;
    required bytes user_token = 4;
}
// 用于登录服务器 验证 返回
// cmd:10021
message login_ack {
}

// 推送下线消息
// cmd:10032
message offline_ntf {
	required int32 type = 1;
}
```
---

# 用户加入场景
用户加入场景时，必须指定是否使用上次的位置`use_last`，如果`use_last=true`，则使用上次离开的位置；如果上次没有登录过该场景，那么使用pb中传入的坐标和角度，如果pb中没有传，则默认值为0；如果`use_last=false`使用pb中传入的坐标和角度，如果pb中没有传，则默认值为0.（注意：pb中3个optional字段需同时都传值才有效，否则当做未传值）

用户加入后，会收到0个或多个`user_list_ntf`的推送消息，该消息推送周围九宫格内的玩家列表数据（包含自己）。

周围九宫格其它玩家会收到`join_ntf`的推送消息，推送新加入的玩家数据。

场景拆分九宫格的宽高在创建场景的API中指定。
```
// 加入场景请求
// cmd:20010
message join_req {
    required bytes space_id = 1;
    optional float pos_x = 2;
    optional float pos_y = 3;
    optional float angle = 4;
    required bool use_last = 5;
    required bytes ex_data = 6;
}
// 加入场景返回
// cmd:20011
message join_ack {
    required int32 width = 1;
    required int32 height = 2;
}
// 加入场景推送
// cmd:20012
message join_ntf {
    required user_data user = 1;
}
// 推送其它格子的玩家列表
// cmd:20013
message user_list_ntf {
    repeated user_data user_list = 1;
}
```
---

# 用户离开场景
用户的tcp链接断开后，会自动离开场景，且其它玩家会收到`leave_ntf`消息
```
// 离开场景请求
// cmd:20020
message leave_req {
    optional bytes space_id = 1;
}
// 离开场景返回
// cmd:20021
message leave_ack {
}
// 离开场景推送
// cmd:20022
message leave_ntf {
    required int32 user_id = 1;
}
```
---

# 用户在场景中移动
用户移动过程中，服务器会判断场景边界，和帧序列。大于边界的值会被修改为等于边界值，即移动不能超出边界；`time`字段指定了每一次移动的帧顺序，服务器处理移动时，小于上次移动真序号的消息会被忽略，即客户端发送移动消息时，`time`必须**大于等于**上次发送的`time`，如果不关心顺序，`time`可一直传相同的值。
```
// 场景移动请求
// cmd:20030
message move_req {
    required float pos_x = 1;
    required float pos_y = 2;
    required float angle = 3;
    required int32 time = 4;
}
// 场景移动返回
// cmd:20031
message move_ack {
    required float pos_x = 1;
    required float pos_y = 2;
    required float angle = 3;
}
// 场景移动推送
// cmd:20032
message move_ntf {
    required int32 user_id = 1;
    required float pos_x = 2;
    required float pos_y = 3;
    required float angle = 4;
}
```
---

# 用户自定义事件广播
广播数据`data`由客户端自定义，该消息可用于附近聊天，玩家动作、技能等。
```
// 场景广播请求
// cmd:20040
message broadcast_req {
    required bytes data = 1;
}
// 场景广播返回
// cmd:20041
message broadcast_ack {
}
// 场景广播推送
// cmd:20042
message broadcast_ntf {
    required int32 user_id = 1;
    required bytes data = 2;
}
```
---

# 查询用户位置
```
// 查询位置请求
// cmd:20050
message query_pos_req {
    required int32 user_id = 1;
}
// 查询位置返回
// cmd:20051
message query_pos_ack {
    required bytes space_id = 1;
    required float pos_x = 2;
    required float pos_y = 3;
    required float angle = 4;
}
```

---



