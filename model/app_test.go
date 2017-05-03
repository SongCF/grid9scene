package model

import (
	"jhqc.com/songcf/scene/pb"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestCreateDeleteApp(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	var app_id = "test_app_id"

	e := CreateApp(app_id, "test_name", "test_key")
	if e != pb.ErrAppAlreadyExist && e != nil {
		t.Errorf("Create app error:%v", e.Desc)
	}
	if !HasApp(app_id) {
		t.Error("1.has app error")
	}
	e = DeleteApp(app_id)
	if e != nil {
		t.Errorf("Delete app error:%v", e.Desc)
	}
	if HasApp(app_id) {
		t.Error("2.has app error")
	}
}
