package model

import (
	"fmt"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestCache(t *testing.T) {
	InitConfTest("../conf.ini")
	InitCache()

	rsp1 := CCPool.Cmd("SET", "scene:test:name", "test set name")
	if rsp1.Err != nil {
		t.Errorf("SET cache error:%v", rsp1.Err)
	}
	rsp2 := CCPool.Cmd("GET", "scene:test:name")
	if rsp2.Err != nil {
		t.Errorf("GET cache error:%v", rsp2.Err)
	}
	rsp2str, err := rsp2.Str()
	if err != nil || rsp2str != "test set name" {
		t.Errorf("GET cache error, str=%v", rsp2str)
	}

	// scene:app_id:grid:space_id:grid_id  ->   uid(set)
	ret := CCPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"), 1)
	if ret.Err != nil {
		t.Errorf("SADD error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"), 2)
	if ret.Err != nil {
		t.Errorf("SADD error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"), 1)
	if ret.Err != nil {
		t.Errorf("SADD error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SMOVE", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"),
		fmt.Sprintf(FORMAT_GRID, "test", "1", "1,1"), 1)
	if ret.Err != nil {
		t.Errorf("SMOVE error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"), 3)
	if ret.Err != nil {
		t.Errorf("SADD error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SREM", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"), 2)
	if ret.Err != nil {
		t.Errorf("SREM error:%v", ret.Err)
	}
	ret = CCPool.Cmd("SMEMBERS", fmt.Sprintf(FORMAT_GRID, "test", "1", "0,0"))
	if ret.Err != nil {
		t.Errorf("SMEMBERS error:%v", ret.Err)
	}

	ret = CCPool.Cmd("SMEMBERS", fmt.Sprintf(FORMAT_GRID, "test", "1", "1,1"))
	if ret.Err != nil {
		t.Errorf("SMEMBERS error:%v", ret.Err)
	}

	// scene:app_id:user:uid  ->  {space_id,grid_id,x,y,angle,exd,node}
	ret = CCPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "test", 1), "space_id", "1", "grid_id", "1")
	if ret.Err != nil {
		t.Errorf("HMSET error:%v", ret.Err)
	}
	ret = CCPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "test", 2), "space_id", "1", "grid_id", "1")
	if ret.Err != nil {
		t.Errorf("HMSET error:%v", ret.Err)
	}
	ret = CCPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "test", 3), "space_id", "1", "grid_id", "1")
	if ret.Err != nil {
		t.Errorf("HMSET error:%v", ret.Err)
	}
	ret = CCPool.Cmd("DEL", fmt.Sprintf(FORMAT_USER, "test", 2))
	if ret.Err != nil {
		t.Errorf("DEL error:%v", ret.Err)
	}

	ret = CCPool.Cmd("HMGET", fmt.Sprintf(FORMAT_USER, "test", 1), "space_id", "grid_id")
	if ret.Err != nil {
		t.Errorf("HMGET error:%v", ret.Err)
	}
}
