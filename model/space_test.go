package model

import (
	"jhqc.com/songcf/scene/pb"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestCreateDeleteSpace(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	var app_id = "test_app_id"
	var space_id = "test_space_id"

	e := CreateSpace(app_id, space_id, float32(10), float32(10))
	if e != pb.ErrSpaceAlreadyExist && e != nil {
		t.Errorf("Create space error:%v", e.Desc)
	}

	w, h, e := GetSpaceInfo(app_id, space_id)
	if e != nil || w != 10 || h != 10 {
		t.Errorf("get space info error:%v, w:%v, h:%v", e.Desc, w, h)
	}

	e = DeleteSpace(app_id, space_id)
	if e != nil {
		t.Errorf("delete space error:%v", e.Desc)
	}
}
