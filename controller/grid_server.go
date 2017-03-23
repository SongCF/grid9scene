package controller

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	. "jhqc.com/songcf/scene/gateway"
)




func Msg2Grid(appId, spaceId, gridId string, m []byte) {
	if grid, ok := AppInfoL[appId].SpaceM[spaceId].GridList[gridId]; ok {
		grid.MsgBox <- m
	} else {
		GlobalWG.Add(1)
		go startGrid(appId, spaceId, gridId, m)
	}
}

func startGrid(appId, spaceId, gridId string, m []byte) {
	defer GlobalWG.Done()

	if grid, ok := AppInfoL[appId].SpaceM[spaceId].GridList[gridId]; ok {
		log.Infof("grid server[%v:%v:%v] already exist.", appId, spaceId, gridId)
		grid.MsgBox <- m
		return
	}
	defer func() {
		log.Debugf("grid server stop. %v:%v:%v", appId, spaceId, gridId)
	}()

	grid := &Grid{
		GridId: gridId,
		MsgBox: make(chan []byte),
	}
	AppInfoL[appId].SpaceM[spaceId].GridList[gridId] = grid
	defer delete(AppInfoL[appId].SpaceM[spaceId].GridList, gridId)

	grid.MsgBox <- m
	// loop
	for {
		select {
		case data, ok := <- grid.MsgBox:
			if !ok {
				return
			}
			log.Infof("handle msg:%v", data)
		// TODO something
		case <- GlobalDie:
			return
		}
	}
}


