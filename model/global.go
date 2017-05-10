package model

import (
	"sync"
	"jhqc.com/songcf/scene/pb"
	"database/sql"
	log "github.com/Sirupsen/logrus"
)


type App struct {
	SsMutex sync.RWMutex
	SessionM map[int32]*Session

	SpaceMutex sync.RWMutex
	SpaceM map[string]*SpaceInfo
}

var appMutex sync.RWMutex
var appPool = map[string]*App{}
func GetAppPool() *map[string]*App {
	return &appPool
}


func getApp(appId string) *App {
	appMutex.RLock()
	app, ok := appPool[appId]
	appMutex.RUnlock() //下面setApp会加锁，所以这里不要用defer
	if !ok {
		app = &App{
			SessionM:make(map[int32]*Session),
			SpaceM:make(map[string]*SpaceInfo),
		}
		setApp(appId, app)
	}
	return app
}
func setApp(appId string, app *App) {
	if app != nil {
		appMutex.Lock()
		defer appMutex.Unlock()
		appPool[appId] = app
	}
}
func DelApp(appId string) {
	appMutex.Lock()
	defer appMutex.Unlock()
	delete(appPool, appId)
}

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


