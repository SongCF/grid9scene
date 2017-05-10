package model

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"jhqc.com/songcf/scene/pb"
	. "jhqc.com/songcf/scene/util"
	"time"
)

const (
	FORMAT_GRID = "scene:%s:grid:%s:%s" //scene:app_id:grid:space_id:grid_id
	FORMAT_USER = "scene:%s:user:%v"    //scene:app_id:user:uid
	NIL         = "nil"
)

var (
	CCPool *pool.Pool
)

func InitCache() {
	addr, err := Conf.Get(SCT_CACHE, "cc_server")
	CheckError(err)
	auth, err := Conf.Get(SCT_CACHE, "cc_auth")
	CheckError(err)
	db, err := Conf.GetInt(SCT_CACHE, "cc_db")
	CheckError(err)
	size, err := Conf.GetInt(SCT_CACHE, "cc_conn_init_size")
	CheckError(err)
	idleTime, err := Conf.GetInt(SCT_CACHE, "cc_conn_idle_time")
	CheckError(err)

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
	CCPool = p

	//clean data
	log.Info("Flush cache...")
	err = CCPool.Cmd("FLUSHDB").Err
	CheckError(err)
}

// except self uid
func GetRoundUidList(appId, spaceId string, gridIdL *[]string, uid int32, conn *redis.Client) ([]int32, *pb.ErrInfo) {
	rg := make([]interface{}, len(*gridIdL))
	for i, gid := range *gridIdL {
		rg[i] = fmt.Sprintf(FORMAT_GRID, appId, spaceId, gid)
	}
	uidUnion := conn.Cmd("SUNION", rg...)
	if uidUnion.Err != nil {
		log.Errorf("GetRoundUidList(user[%v:%v]) Cache Cmd(SUNION) error(%v)", appId, uid, uidUnion.Err)
		return nil, pb.ErrServerBusy
	}
	unionRespL, err := uidUnion.Array()
	if err != nil {
		log.Errorf("GetRoundUidList(user[%v:%v]) get uidUnion array error(%v)", appId, uid, err)
		return nil, pb.ErrServerBusy
	}
	uidL := []int32{}
	for _, resp := range unionRespL {
		tmpUid, err := resp.Int()
		if err != nil {
			log.Errorf("GetRoundUidList(user[%v:%v]) parse uid int error(%v)", appId, uid, err)
			return nil, pb.ErrServerBusy
		}
		tu := int32(tmpUid)
		if uid != tu { // except self
			uidL = append(uidL, tu)
		}
	}
	return uidL, nil
}

func GetUserInfo(appId string, uid int32, conn *redis.Client) (*UserInfo, *pb.ErrInfo) {
	resp := conn.Cmd("HMGET", fmt.Sprintf(FORMAT_USER, appId, uid),
		StrSpaceId, StrGridId, StrX, StrY, StrAngle, StrMoveTime, StrExData)
	if resp.Err != nil {
		log.Errorf("GetUserInfo user(%v:%v) data HMGET error: %v", appId, uid, resp.Err)
		return nil, pb.ErrServerBusy
	}
	l, err := resp.Array()
	if err != nil || len(l) != 7 {
		log.Errorf("GetUserInfo user(%v:%v) data parse array error: %v", appId, uid, resp.Err)
		return nil, pb.ErrServerBusy
	}
	if l[0].IsType(redis.Nil) {
		if !HasApp(appId) {
			return nil, pb.ErrAppNotExist
		}
		return nil, pb.ErrUserOffline
	}
	spaceId, err0 := l[0].Str()
	gridId, err1 := l[1].Str()
	x, err2 := l[2].Float64()
	y, err3 := l[3].Float64()
	angle, err4 := l[4].Float64()
	moveTime, err5 := l[5].Int()
	exd, err6 := l[6].Str()
	if err0 != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		log.Errorf("GetUserInfo(user[%v:%v]) parse userInfo error(%v,%v,%v,%v,%v,%v,%v), resp:%v",
			appId, uid, err0, err1, err2, err3, err4, err5, err6, l)
		return nil, pb.ErrServerBusy
	}
	return &UserInfo{
		SpaceId:  spaceId,
		GridId:   gridId,
		PosX:     float32(x),
		PosY:     float32(y),
		Angle:    float32(angle),
		MoveTime: int32(moveTime),
		ExData:   []byte(exd),
	}, nil
}

func ResetUserInfo(appId string, uid int32, conn *redis.Client) *pb.ErrInfo {
	// UserInfo
	err := conn.Cmd("HMSET", fmt.Sprintf(FORMAT_USER, appId, uid),
		StrSpaceId, NIL, StrGridId, NIL,
		StrX, DEFAULT_POS_X, StrY, DEFAULT_POS_Y, StrAngle, DEFAULT_ANGLE,
		StrMoveTime, 0, StrExData, "").Err
	if err != nil {
		return pb.ErrServerBusy
	}
	return nil
}
