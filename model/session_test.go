package model

import (
	"testing"
)

func TestSetGetSession(t *testing.T) {
	var app_id = "test_app_id"
	var uid = int32(1)

	var s Session
	s.AppId = app_id
	s.Uid = uid
	if !s.HasLogin() {
		t.Error("user has login error")
	}
	//clean
	s.Clean()
	if s.AppId != "" || s.Uid != 0 {
		t.Error("clean session error")
	}
}
