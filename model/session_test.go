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
	SetSession(app_id, uid, &s)
	if &s != GetSession(app_id, uid) {
		t.Error("get session error")
	}
}
