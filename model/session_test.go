package model

import (
	"jhqc.com/songcf/scene/_test"
	"testing"
)

func TestSetGetSession(t *testing.T) {
	var s Session
	s.AppId = _test.T_APP_ID
	s.Uid = _test.T_USER_ID
	if !s.HasLogin() {
		t.Error("user has login error")
	}
	SetSession(_test.T_APP_ID, _test.T_USER_ID, &s)
	if &s != GetSession(_test.T_APP_ID, _test.T_USER_ID) {
		t.Error("get session error")
	}
}
