package _test

import (
	"encoding/json"
	"jhqc.com/songcf/scene/util"
	"runtime/debug"
	"testing"
)

func StartClient(t *testing.T) {
	//从zk获取服务地址
	httpSrvL := util.GetServices("http")
	if len(httpSrvL) != 1 && len(httpSrvL[0]) != 0 {
		t.Fatalf("get http server error, ret:%v", httpSrvL)
	}
	zkHttp := &util.ZkData{}
	err := json.Unmarshal(httpSrvL[0], zkHttp)
	check(err, t)

	tcpSrvL := util.GetServices("tcp")
	if len(tcpSrvL) != 1 && len(tcpSrvL[0]) != 0 {
		t.Fatalf("get tcp server error, ret:%v", tcpSrvL)
	}
	zkTcp := &util.ZkData{}
	err = json.Unmarshal(tcpSrvL[0], zkTcp)
	check(err, t)

	//创建测试 app_id space_id
	testAllAPI(zkHttp.URI, t)
	//建立tcp连接
	testAllMsg(zkHttp.URI, zkTcp.URI, t)
}

func check(err error, t *testing.T) {
	if err != nil {
		debug.PrintStack()
		t.Fatal(err)
	}
}

func assert(b bool, desc string, t *testing.T) {
	if !b {
		debug.PrintStack()
		t.Fatal(desc)
	}
}
