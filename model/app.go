package model

import (
	log "github.com/Sirupsen/logrus"
	"jhqc.com/songcf/scene/pb"
)


func CreateApp(appId, name, key string) *pb.ErrInfo {
	_, err := DB.Exec("INSERT INTO ?(app_id,name,private_key) values(?,?,?);",
		TBL_APP, appId, name, key)
	if err != nil {
		log.Errorf("create app failed, appid=%v", appId)
		return pb.ErrQueryDBError
	}
	AppInfoL[appId] = &AppInfo{}
	return nil
}

func DeleteApp(appId string) *pb.ErrInfo {
	_, err := DB.Exec("DELETE FROM ? WHERE app_id=?;", TBL_APP, appId)
	if err != nil {
		log.Errorf("delete app failed, appid=%v", appId)
		return pb.ErrQueryDBError
	}
	//TODO stop grid server

	//clean cache
	delete(AppInfoL, appId)
	return nil
}

func HasApp(appId string) bool {
	if _,ok := AppInfoL[appId]; ok {
		return true
	}
	return false
}
