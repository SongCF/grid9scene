package controller

import (
	. "jhqc.com/songcf/scene/model"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	"sync"
	"jhqc.com/songcf/scene/pb"
)



func Msg2AppWait(appId string, imsg int32, spaceId string) {
	app := GetApp(appId)
	if app == nil {
		app = StartApp(appId)
	}
	if app == nil {
		log.Infof("not found app[%v], ignore msg[%v]", appId, imsg)
		return
	}
	cb := make(chan struct{})
	msg := &InnerMsg{
		Id: imsg,
		Cb: cb,
		SpaceId: spaceId,
	}
	app.PostMsg(msg)
	<- cb //waiting close cb
}


var mutex *sync.Mutex = new(sync.Mutex)
func StartApp(appId string) *App {
	mutex.Lock() //http可能会多个一起请求，避免创建同一app多次
	ch := make(chan struct{})
	go appServe(appId, ch)
	<- ch
	mutex.Unlock()
	app := GetApp(appId)
	return app
}

func appServe(appId string, ch chan struct{}) {
	alreadyApp := GetApp(appId)
	if alreadyApp != nil {
		log.Infof("app server[%v] already exist.", appId)
		close(ch) //started.
		return
	}

	// query app info from db
	raw := DB.QueryRow("SELECT app_id FROM ? WHERE app_id=?;",
		TBL_APP, appId)
	var q string
	err := raw.Scan(&q) // if empty, err = sql.ErrNoRows
	if err != nil {
		log.Errorf("query db error = %v\n", err)
		close(ch)
		return
	}

	//cache
	app := &App{
		SpaceM: make(map[string]*Space),
		SessionM: make(map[int32]*Session),
		MsgBox: make(chan *InnerMsg),
	}
	AppL[appId] = app
	defer func() {
		//delete cache
		delete(AppL, appId)
		log.Infof("app server stop. %v", appId)
	}()

	close(ch) //started.

	// loop
	for {
		select {
		case data, ok := <- app.MsgBox:
			if !ok {
				return
			}
			log.Infof("handle app msg:%v", data)
			switch data.Id {
			case IMSG_START_SPACE:
				StartSpace(appId, data.SpaceId)
				close(data.Cb)
			default:
				log.Infof("app_server handle unknow msg:%v", data.Id)
			}
		case <- GlobalDie:
			return
		}
	}
}





func CreateApp(appId, name, key string) *pb.ErrInfo {
	//already exist
	app := GetApp(appId)
	if app != nil {
		return nil
	}
	//db
	_, err := DB.Exec("INSERT INTO ?(app_id,name,private_key) values(?,?,?);",
		TBL_APP, appId, name, key)
	if err != nil {
		log.Errorf("create app failed, appid=%v", appId)
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
	tx.Exec("DELETE FROM ? WHERE app_id=?;", TBL_APP, appId)
	tx.Exec("DELETE FROM ? WHERE app_id=?;", TBL_SPACE, appId)
	tx.Exec("DELETE FROM ? WHERE app_id=?;", TBL_LAST_SPACE, appId)
	tx.Exec("DELETE FROM ? WHERE app_id=?;", TBL_LAST_POS, appId)
	err = tx.Commit()
	if err != nil {
		log.Errorln("delete app failed, db commit failed")
		return pb.ErrQueryDBError
	}

	app := GetApp(appId)
	if app != nil {
		//delete user session
		for _,s := range app.SessionM {
			s.Close()
		}
		//delete cache, close server
		for _,space := range app.SpaceM {
			for _,grid := range space.GridM {
				grid.Close()
			}
			space.Close()
		}
	}
	app.Close()
	return nil
}

