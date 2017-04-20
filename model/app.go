package model

import (
	log "github.com/Sirupsen/logrus"
	"jhqc.com/songcf/scene/pb"
)

func HasApp(appId string) bool {
	//db
	raw := DB.QueryRow("SELECT app_id from app WHERE app_id=?;", appId)
	var tmpAppId string
	err := raw.Scan(&tmpAppId) // if empty, err = sql.ErrNoRows
	if err != nil || tmpAppId != appId {
		return false
	} else {
		return true
	}
}

func CreateApp(appId, name, key string) *pb.ErrInfo {
	//db
	_, err := DB.Exec("INSERT INTO app(app_id,name,private_key) values(?,?,?);",
		appId, name, key)
	if IsDuplicate(err) {
		return pb.ErrAppAlreadyExist
	}
	if err != nil {
		log.Errorf("create app failed, appid=%v, err=%v", appId, err)
		return pb.ErrQueryDBError
	}
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
	return nil
}
