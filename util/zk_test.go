package util

import "testing"

func TestRegister(t *testing.T) {
	Register()
	if ZkConn == nil {
		t.Error("Connect ZK Error")
	}
	tcpSrv := GetServices("tcp")
	if len(tcpSrv) < 1 {
		t.Errorf("Register servers not found, tcpSrv:%v", tcpSrv)
	}
	if string(tcpSrv[0]) == "" {
		t.Errorf("Register server is none, tcpSrv[0]:%v", tcpSrv)
	}
	// unregister
	if err := ZkConn.Delete(zkTcpKey, -1); err != nil {
		t.Errorf("unregister tcp server error:%v", err)
	}
	if err := ZkConn.Delete(zkHttpKey, -1); err != nil {
		t.Errorf("unregister http server error:%v", err)
	}
}
