package controller

import (
	. "jhqc.com/songcf/scene/model"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	"jhqc.com/songcf/scene/pb"
	"database/sql"
)

func Msg2SpaceWait(appId, spaceId string, imsg int32, gridId string) {
	space := GetSpace(appId, spaceId)
	if space == nil {
		Msg2AppWait(appId, IMSG_START_SPACE, spaceId)
		space = GetSpace(appId, spaceId)
	}
	if space == nil {
		log.Infof("start space_server[%v:%v] failed, ignore msg[%v]", appId, spaceId, imsg)
		return
	}
	cb := make(chan struct{})
	msg := &InnerMsg{
		Id: imsg,
		Cb: cb,
		GridId: gridId,
	}
	space.PostMsg(msg)
	<- cb //waiting close cb
}


// unsafe multi goruntine
func StartSpace(appId, spaceId string) *Space {
	ch := make(chan struct{})
	go spaceServe(appId, spaceId, ch)
	<- ch
	space := GetSpace(appId, spaceId)
	return space
}



func spaceServe(appId, spaceId string, ch chan struct{}) {
	// check already started.
	alreadySpace := GetSpace(appId, spaceId)
	if alreadySpace != nil {
		log.Infof("space server[%v:%v] already exist.", appId, spaceId)
		close(ch) //started.
		return
	}

	app := GetApp(appId)
	if app == nil {
		log.Errorf("not found app = %v \n", appId)
		close(ch)
		return
	}

	// query space info from db
	raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;", appId, spaceId)
	var w, h float32
	err := raw.Scan(&w, &h) // if empty, err = sql.ErrNoRows
	if err == sql.ErrNoRows {
		log.Infof("Space(%v:%v) doesn't exist", appId, spaceId)
		close(ch) //started.
		return
	}
	if err != nil {
		log.Errorf("select grid w h error(%v:%v) = %v\n", appId, spaceId, err)
		close(ch) //started.
		return
	}

	//cache
	space := &Space{
		SpaceId: spaceId,
		GridWidth: w,
		GridHeight: h,
		GridM: make(map[string]*Grid),
		MsgBox: make(chan *InnerMsg),
	}
	app.SpaceM[spaceId] = space
	defer func() {
		//delete cache
		if app, ok := AppL[appId]; ok {
			delete(app.SpaceM, spaceId)
		}
		log.Infof("space server stop. %v:%v", appId, spaceId)
	}()

	close(ch) //started.

	// loop
	for {
		select {
		case data, ok := <- space.MsgBox:
			if !ok {
				return
			}
			log.Infof("handle app msg:%v", data)
			switch data.Id {
			case IMSG_START_GRID:
				StartGrid(appId, spaceId, data.GridId)
				close(data.Cb)
			default:
				log.Infof("space_server handle unknow msg")
			}
		case <- GlobalDie:
			return
		}
	}
}



func CreateSpace(appId, spaceId string, gridWidth, gridHeight float32) *pb.ErrInfo {
	//already exist
	alreadySpace := GetSpace(appId, spaceId)
	if alreadySpace != nil {
		return nil
	}
	//check app
	app := GetApp(appId)
	if app == nil {
		return pb.ErrAppNotExist
	}
	//db
	_, err := DB.Exec("INSERT INTO space(app_id,space_id,grid_width,grid_height) values(?,?,?,?);",
		appId, spaceId, gridWidth, gridHeight)
	if err != nil && !IsDuplicate(err) {
		log.Errorf("create space(%v:%v) failed, err=%v", appId, spaceId, err)
		return pb.ErrQueryDBError
	}
	//start space_server
	Msg2AppWait(appId, IMSG_START_SPACE, spaceId)
	return nil
}

func DeleteSpace(appId, spaceId string) *pb.ErrInfo {
	//delete db
	tx, err := DB.Begin()
	if err != nil {
		log.Errorln("delete space failed, db begin failed")
		return pb.ErrQueryDBError
	}
	DB.Exec("DELETE FROM space WHERE app_id=? and space_id=?;", appId, spaceId)
	tx.Exec("DELETE FROM last_space WHERE app_id=? and space_id=?;", appId, spaceId)
	tx.Exec("DELETE FROM last_pos WHERE app_id=? and space_id=?;", appId, spaceId)
	err = tx.Commit()
	if err != nil {
		log.Errorln("delete space failed, db commit failed")
		return pb.ErrQueryDBError
	}
	//delete cache, close server
	space := GetSpace(appId, spaceId)
	if space != nil {
		for _,grid := range space.GridM {
			grid.Close()
		}
	}
	space.Close()
	return nil
}

