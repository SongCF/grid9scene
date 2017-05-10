package model

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"jhqc.com/songcf/scene/pb"
	"sync"
	"time"
)

type App struct {
	SsMutex  sync.RWMutex `json:"-"`
	SessionM map[int32]*Session

	SpaceMutex sync.RWMutex `json:"-"`
	SpaceM     map[string]*SpaceInfo
}

var appPoolMutex sync.RWMutex
var appPool = map[string]*App{}

//=====================================================
//app
//=====================================================

func getApp(appId string) *App {
	appPoolMutex.RLock()
	app, ok := appPool[appId]
	appPoolMutex.RUnlock() //下面setApp会加锁，所以这里不要用defer
	if !ok {
		app = &App{
			SessionM: make(map[int32]*Session),
			SpaceM:   make(map[string]*SpaceInfo),
		}
		setApp(appId, app)
	}
	return app
}
func setApp(appId string, app *App) {
	if app != nil {
		appPoolMutex.Lock()
		defer appPoolMutex.Unlock()
		appPool[appId] = app
	}
}
func DelApp(appId string) {
	appPoolMutex.Lock()
	defer appPoolMutex.Unlock()
	delete(appPool, appId)
}

//=====================================================
//session
//=====================================================

func SetSession(appId string, uid int32, s *Session) {
	app := getApp(appId)
	app.SsMutex.Lock()
	defer app.SsMutex.Unlock()
	app.SessionM[uid] = s
}
func GetSession(appId string, uid int32) *Session {
	app := getApp(appId)
	app.SsMutex.RLock()
	defer app.SsMutex.RUnlock()
	s, ok := app.SessionM[uid]
	if ok {
		return s
	}
	return nil
}
func DelSession(appId string, uid int32) {
	app := getApp(appId)
	app.SsMutex.Lock()
	defer app.SsMutex.Unlock()
	delete(app.SessionM, uid)
}

//调试用:获取关心的session信息
func GetAllSession() *map[string]string {
	m := make(map[string]string)
	appPoolMutex.RLock()
	defer appPoolMutex.RUnlock()
	for tAppId, tApp := range appPool {
		tApp.SsMutex.RLock()
		for tUid, tSs := range tApp.SessionM {
			m[fmt.Sprintf("%v:%v", tAppId, tUid)] =
				fmt.Sprintf("pk_count:%v, ip:%v, conn_time:%v", tSs.PacketCount, tSs.IP, tSs.ConnectTime)
		}
		tApp.SsMutex.RUnlock()
	}
	return &m
}

//调试用：获取所有session的数量，数据包总计，链接总时间毫秒
func GetAllSessionMsgAvg() (int, int, int) {
	pkCount := 0
	timeAll := time.Time{}
	ssNum := 0
	appPoolMutex.RLock()
	defer appPoolMutex.RUnlock()
	for _, tApp := range appPool {
		tApp.SsMutex.RLock()
		for _, tSs := range tApp.SessionM {
			pkCount += int(tSs.PacketCount)
			timeAll = timeAll.Add(time.Now().Sub(tSs.ConnectTime))
			ssNum++
		}
		tApp.SsMutex.RUnlock()
	}
	return ssNum, pkCount, timeAll.Nanosecond() / 1000000
}

//=====================================================
//space
//=====================================================

func setSpace(appId, spaceId string, s *SpaceInfo) {
	app := getApp(appId)
	app.SpaceMutex.Lock()
	defer app.SpaceMutex.Unlock()
	app.SpaceM[spaceId] = s
}
func getSpace(appId, spaceId string) *SpaceInfo {
	app := getApp(appId)
	app.SpaceMutex.Lock()
	defer app.SpaceMutex.Unlock()
	s, ok := app.SpaceM[spaceId]
	if ok {
		return s
	}
	return nil
}
func DelSpace(appId, spaceId string) {
	app := getApp(appId)
	app.SpaceMutex.Lock()
	defer app.SpaceMutex.Unlock()
	delete(app.SpaceM, spaceId)
}

func GetSpaceInfo(appId, spaceId string) (gridWidth, gridHeight float32, e *pb.ErrInfo) {
	if sp := getSpace(appId, spaceId); sp != nil {
		gridWidth = sp.GridWidth
		gridHeight = sp.GridHeight
		return
	}
	// query space info from db
	raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;", appId, spaceId)
	err := raw.Scan(&gridWidth, &gridHeight) // if empty, err = sql.ErrNoRows
	if err == sql.ErrNoRows {
		log.Infof("Get Space(%v:%v) doesn't exist", appId, spaceId)
		e = pb.ErrSpaceNotExist
		return
	}
	if err != nil {
		log.Errorf("select grid w h error(%v:%v) = %v\n", appId, spaceId, err)
		e = pb.ErrQueryDBError
		return
	}
	s := &SpaceInfo{
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
	setSpace(appId, spaceId, s)
	return
}
