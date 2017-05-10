package main

import (
	"encoding/json"
	"jhqc.com/songcf/scene/_test"
	"jhqc.com/songcf/scene/util"
	"testing"
	"time"
)

func TestScene(t *testing.T) {
	// start server
	go main()
	// waiting started
	select {
	case <-server_started:
	case <-time.After(time.Second * 60):
		t.Fatal("start server timeout!")
	}
	// start test client
	startTest(t)
}

func startTest(t *testing.T) {
	//从zk获取服务地址
	httpSrvL := util.GetServices("http")
	if len(httpSrvL) != 1 && len(httpSrvL[0]) != 0 {
		t.Fatalf("get http server error, ret:%v", httpSrvL)
	}
	zkHttp := &util.ZkData{}
	err := json.Unmarshal(httpSrvL[0], zkHttp)
	if err != nil {
		t.Fatal("unmarshal zkhttp json failed:" + err.Error())
	}

	tcpSrvL := util.GetServices("tcp")
	if len(tcpSrvL) != 1 && len(tcpSrvL[0]) != 0 {
		t.Fatalf("get tcp server error, ret:%v", tcpSrvL)
	}
	zkTcp := &util.ZkData{}
	err = json.Unmarshal(tcpSrvL[0], zkTcp)
	if err != nil {
		t.Fatal("unmarshal zktcp json failed:" + err.Error())
	}

	//创建测试 app_id space_id
	_test.TestAllAPI(zkHttp.URI, t)
	//建立tcp连接
	_test.TestAllPB(zkHttp.URI, zkTcp.URI, t)
}
