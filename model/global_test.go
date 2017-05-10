package model

import (
	"testing"
)

func TestApp(t *testing.T) {
	var app_id = "test_app_id"
	DelApp(app_id)
	app := getApp(app_id)
	if app == nil {
		t.Fatal("get app nil")
	}
	DelApp(app_id)
}

func TestSession(t *testing.T) {
	var app_id = "test_app_id"
	var uid = int32(1)
	var s Session
	s.AppId = app_id
	s.Uid = uid

	SetSession(app_id, uid, &s)
	if &s != GetSession(app_id, uid) {
		t.Fatal("set get session error")
	}
	DelSession(app_id, uid)
	if GetSession(app_id, uid) != nil {
		t.Fatal("del get session error")
	}
}

func TestSpace(t *testing.T) {
	var app_id = "test_app_id"
	var space_id = "test_space_id"
	var x, y = float32(10), float32(10)
	if getSpace(app_id, space_id) != nil {
		t.Fatal("get nil space error")
	}
	s := &SpaceInfo{GridWidth: x, GridHeight: y}
	setSpace(app_id, space_id, s)
	if s != getSpace(app_id, space_id) {
		t.Fatal("set get space error")
	}
	w, h, e := GetSpaceInfo(app_id, space_id)
	if e != nil || w != x || h != y {
		t.Fatal("get space info error")
	}
	DelSpace(app_id, space_id)
	if getSpace(app_id, space_id) != nil {
		t.Fatal("del get nil space error")
	}
}
