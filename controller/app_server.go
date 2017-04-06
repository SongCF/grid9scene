package controller

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/pb"
	"jhqc.com/songcf/scene/util"
	"sync"
)

func Msg2AppWait(appId string, imsg int32, spaceId string) {
	app := GetApp(appId)
	if app == nil {
		app = StartApp(appId)
	}
	if app == nil {
		log.Errorf("not found app[%v], ignore msg[%v]", appId, imsg)
		return
	}
	cb := make(chan struct{})
	msg := &InnerMsg{
		Id:      imsg,
		Cb:      cb,
		SpaceId: spaceId,
	}
	app.PostMsg(msg)
	<-cb //waiting close cb
}

var mutex *sync.Mutex = new(sync.Mutex)

func StartApp(appId string) *App {
	mutex.Lock() //http可能会多个一起请求，避免创建同一app多次
	ch := make(chan struct{})
	go appServe(appId, ch)
	<-ch
	mutex.Unlock()
	app := GetApp(appId)
	return app
}

func appServe(appId string, ch chan struct{}) {
	defer util.RecoverPanic()
	alreadyApp := GetApp(appId)
	if alreadyApp != nil {
		log.Infof("app server[%v] already exist.", appId)
		close(ch) //started.
		return
	}
	// query app info from db
	raw := DB.QueryRow("SELECT app_id FROM app WHERE app_id=?;", appId)
	var q string
	err := raw.Scan(&q) // if empty, err = sql.ErrNoRows
	if err == sql.ErrNoRows {
		log.Infof("App(%v) doesn't exist.", appId)
		close(ch)
		return
	}
	if err != nil {
		log.Errorf("select appid error(%v) = %v\n", appId, err)
		close(ch)
		return
	}

	//cache
	app := &App{
		SpaceM:   make(map[string]*Space),
		SessionM: make(map[int32]*Session),
		MsgBox:   make(chan *InnerMsg),
		Die:      make(chan struct{}),
	}
	AppL[appId] = app
	defer func() {
		//delete cache
		if app != nil {
			//delete user session
			for _, s := range app.SessionM {
				s.Close()
			}
			//delete cache, close server
			for _, space := range app.SpaceM {
				for _, grid := range space.GridM {
					grid.Close()
				}
				space.Close()
			}
		}
		delete(AppL, appId)
		log.Infof("app server stop. %v", appId)
	}()

	close(ch) //started.

	// loop
	for {
		select {
		case data, ok := <-app.MsgBox:
			if !ok {
				return
			}
			handleAppMsg(appId, data)
		case <-app.Die:
			return
		case <-GlobalDie:
			return
		}
	}
}

func handleAppMsg(appId string, m *InnerMsg) {
	defer util.RecoverPanic()
	log.Infof("handle app msg:%v", m.Id)
	switch m.Id {
	case IMSG_START_SPACE:
		StartSpace(appId, m.SpaceId)
		close(m.Cb)
	default:
		log.Infof("app_server handle unknow msg:%v", m.Id)
	}
}

func CreateApp(appId, name, key string) *pb.ErrInfo {
	//already exist
	app := GetApp(appId)
	if app != nil {
		return nil
	}
	//db
	_, err := DB.Exec("INSERT INTO app(app_id,name,private_key) values(?,?,?);",
		appId, name, key)
	if err != nil && !IsDuplicate(err) {
		log.Errorf("create app failed, appid=%v, err=%v", appId, err)
		return pb.ErrQueryDBError
	}
	// start app_server
	StartApp(appId)
	return nil
}

func DeleteApp(appId string) *pb.ErrInfo {
	//delete db
	tx, err := DB.Begin()
	if err != nil {
		log.Errorln("delete app failed, db begin failed")
		return pb.ErrQueryDBError
	}
	tx.Exec("DELETE FROM app WHERE app_id=?;", appId)
	tx.Exec("DELETE FROM space WHERE app_id=?;", appId)
	tx.Exec("DELETE FROM last_space WHERE app_id=?;", appId)
	tx.Exec("DELETE FROM last_pos WHERE app_id=?;", appId)
	err = tx.Commit()
	if err != nil {
		log.Errorln("delete app failed, db commit failed")
		return pb.ErrQueryDBError
	}

	app := GetApp(appId)
	if app != nil {
		//delete user session
		for _, s := range app.SessionM {
			s.Close()
		}
		//delete cache, close server
		for _, space := range app.SpaceM {
			for _, grid := range space.GridM {
				grid.Close()
			}
			space.Close()
		}
	}
	app.Close()
	return nil
}
