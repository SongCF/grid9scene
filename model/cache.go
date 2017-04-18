package model

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/pool"
	. "jhqc.com/songcf/scene/util"
	"time"
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
	rsp := ccPool.Cmd("GET", "scene:test:name")
	log.Infof("ret : %v", rsp)
}
