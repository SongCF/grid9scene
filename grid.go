package main

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/types"
)





func msg2Grid(appId, spaceId, gridId string, m []byte) {
	if grid, ok := AppList[appId].SpaceList[spaceId].GridList[gridId]; ok {
		grid.ChanMsg <- m
	} else {
		globalWG.Add(1)
		go startGrid(appId, spaceId, gridId, m)
	}
}

func startGrid(appId, spaceId, gridId string, m []byte) {
	defer globalWG.Done()

	if grid, ok := AppList[appId].SpaceList[spaceId].GridList[gridId]; ok {
		log.Infof("grid server[%v:%v:%v] already exist.", appId, spaceId, gridId)
		grid.ChanMsg <- m
		return
	}
	defer func() {
		log.Debugf("grid server stop. %v:%v:%v", appId, spaceId, gridId)
	}()

	grid := &Grid{
		GridId: gridId,
		ChanMsg: make(chan []byte),
	}
	AppList[appId].SpaceList[spaceId].GridList[gridId] = grid
	defer delete(AppList[appId].SpaceList[spaceId].GridList, gridId)

	grid.ChanMsg <- m
	// loop
	for {
		select {
		case data, ok := <- grid.ChanMsg:
			if !ok {
				return
			}
			log.Infof("handle msg:%v", data)
			// TODO something
		case <- globalDie:
			return
		}
	}
}


