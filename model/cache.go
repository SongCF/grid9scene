package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/pool"
	. "jhqc.com/songcf/scene/util"
	"time"
	"fmt"
)


const (
	FORMAT_GRID = "scene:%s:grid:%s:%s"   //scene:app_id:grid:space_id:grid_id
	FORMAT_USER = "scene:%s:user:%v" //scene:app_id:user:uid
)

// mem cache

var (
	AppL = make(map[string]*App)
	ccPool *pool.Pool
)


func InitCache() {
	addr := Conf.Get(SCT_CACHE, "cc_server")
	if addr == "" {
		panic("get cache server addr error")
	}
	auth := Conf.Get(SCT_CACHE, "cc_auth")
	if auth == "" {
		panic("get cache auth error")
	}
	db := Conf.GetInt(SCT_CACHE, "cc_db", -1)
	if db == -1 {
		panic("get cache db error")
	}
	size := Conf.GetInt(SCT_CACHE, "cc_max_open_conn", 2000)
	idleTime := Conf.GetInt(SCT_CACHE, "cc_conn_idle_time", 60)
	log.Infof("redis addr:%v, auth:%v, size:%v, idletime:%v", addr, auth, size, idleTime)

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			log.Error("redis connect err:", err)
			return nil, err
		}
		if err = client.Cmd("AUTH", auth).Err; err != nil {
			client.Close()
			log.Error("redis auth err:", err)
			return nil, err
		}
		if err = client.Cmd("SELECT", db).Err; err != nil {
			client.Close()
			return nil, err
		}
		go func() {
			for {
				client.Cmd("PING")
				time.Sleep(time.Duration(idleTime) * time.Second)
			}
		}()
		return client, nil
	}
	p, err := pool.NewCustom("tcp", addr, size, df)
	CheckError(err)
	ccPool = p

	//TODO delete
	test()
}


func test() {
	log.Debug("cache test...")
	rsp1 := ccPool.Cmd("SET", "scene:test:name", "test set name")
	rsp2 := ccPool.Cmd("GET", "scene:test:name")
	log.Infof("set ret:%v, get ret:%v", rsp1, rsp2)

	// scene:app_id:grid:space_id:grid_id  ->   uid(set)
	ret := ccPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), 1)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), 2)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), 1)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SMOVE", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), fmt.Sprintf(FORMAT_GRID, "1", "1", "1,1"), 1)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SADD", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), 3)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SREM", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"), 2)
	CheckError(ret.Err)
	ret = ccPool.Cmd("SMEMBERS", fmt.Sprintf(FORMAT_GRID, "1", "1", "0,0"))
	CheckError(ret.Err)
	log.Infof("set(0,0) mem:%v", ret)
	ret = ccPool.Cmd("SMEMBERS", fmt.Sprintf(FORMAT_GRID, "1", "1", "1,1"))
	CheckError(ret.Err)
	log.Infof("set(1,1) mem:%v", ret)

	// scene:app_id:user:uid  ->  {space_id,grid_id,x,y,angle,exd,node}
	ret = ccPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "1", 1), "space_id", "1", "grid_id", "1")
	CheckError(ret.Err)
	ret = ccPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "1", 2), "space_id", "1", "grid_id", "1")
	CheckError(ret.Err)
	ret = ccPool.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, "1", 3), "space_id", "1", "grid_id", "1")
	CheckError(ret.Err)
	ret = ccPool.Cmd("DEL", fmt.Sprintf(FORMAT_USER, "1", 2))
	CheckError(ret.Err)

	// pipeline
	ret = ccPool.Cmd("HMGET", fmt.Sprintf(FORMAT_USER, "1", 1), "space_id", "grid_id")
	CheckError(ret.Err)
	log.Infof("HASH(1,1):%v", ret)
}
