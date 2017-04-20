package model

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"jhqc.com/songcf/scene/pb"
)

type spaceInfo struct {
	GridWidth  float32
	GridHeight float32
}

var spaceM = map[string]*spaceInfo{}

func CreateSpace(appId, spaceId string, gridWidth, gridHeight float32) *pb.ErrInfo {
	//already exist
	_, _, e := GetSpaceInfo(appId, spaceId)
	if e == nil { //成功获取到了，存在该space
		return pb.ErrSpaceAlreadyExist
	} else if e != pb.ErrSpaceNotExist {
		return e
	}
	//check app
	if !HasApp(appId) {
		return pb.ErrAppNotExist
	}
	//db
	_, err := DB.Exec("INSERT INTO space(app_id,space_id,grid_width,grid_height) values(?,?,?,?);",
		appId, spaceId, gridWidth, gridHeight)
	if IsDuplicate(err) {
		return pb.ErrSpaceAlreadyExist
	}
	if err != nil {
		log.Errorf("create space(%v:%v) failed, err=%v", appId, spaceId, err)
		return pb.ErrQueryDBError
	}
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
	return nil
}

func GetSpaceInfo(appId, spaceId string) (gridWidth, gridHeight float32, e *pb.ErrInfo) {
	if info, ok := spaceM[fmt.Sprintf("%v:%v", appId, spaceId)]; ok {
		gridWidth = info.GridWidth
		gridHeight = info.GridHeight
		return
	}
	// query space info from db
	raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;", appId, spaceId)
	err := raw.Scan(&gridWidth, &gridHeight) // if empty, err = sql.ErrNoRows
	if err == sql.ErrNoRows {
		log.Infof("Space(%v:%v) doesn't exist", appId, spaceId)
		e = pb.ErrSpaceNotExist
		return
	}
	if err != nil {
		log.Errorf("select grid w h error(%v:%v) = %v\n", appId, spaceId, err)
		e = pb.ErrQueryDBError
		return
	}
	spaceM[fmt.Sprintf("%v:%v", appId, spaceId)] = &spaceInfo{
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
	return
}
