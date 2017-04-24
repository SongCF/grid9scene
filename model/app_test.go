package model

import (
	"jhqc.com/songcf/scene/_test"
	"jhqc.com/songcf/scene/pb"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestCreateDeleteApp(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	e := CreateApp(_test.T_APP_ID, "test_name", "test_key")
	if e != pb.ErrAppAlreadyExist && e != nil {
		t.Errorf("Create app error:%v", e.Desc)
	}
	if !HasApp(_test.T_APP_ID) {
		t.Error("1.has app error")
	}
	e = DeleteApp(_test.T_APP_ID)
	if e != nil {
		t.Errorf("Delete app error:%v", e.Desc)
	}
	if HasApp(_test.T_APP_ID) {
		t.Error("2.has app error")
	}
}
