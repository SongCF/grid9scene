package model

import (
	"jhqc.com/songcf/scene/_test"
	"jhqc.com/songcf/scene/pb"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestCreateDeleteSpace(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	e := CreateSpace(_test.T_APP_ID, _test.T_SPACE_ID, float32(10), float32(10))
	if e != pb.ErrSpaceAlreadyExist && e != nil {
		t.Errorf("Create space error:%v", e.Desc)
	}

	w, h, e := GetSpaceInfo(_test.T_APP_ID, _test.T_SPACE_ID)
	if e != nil || w != 10 || h != 10 {
		t.Errorf("get space info error:%v, w:%v, h:%v", e.Desc, w, h)
	}

	e = DeleteSpace(_test.T_APP_ID, _test.T_SPACE_ID)
	if e != nil {
		t.Errorf("delete space error:%v", e.Desc)
	}
}
