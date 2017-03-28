package model


import (
	log "github.com/Sirupsen/logrus"
)




// mem cache
type AppInfo struct {
	SpaceM map[string]*Space      // spaceId : Space
	SessionM map[int32]*Session // uid : Session
}
var (
	AppInfoL = make(map[string]*AppInfo)
)



func GetGrid(appId, spaceId, gridId string) *Grid {
	if app, ok := AppInfoL[appId]; ok {
		if space, ok := app.SpaceM[spaceId]; ok {
			if grid, ok := space.GridM[gridId]; ok {
				return grid
			}
		}
	}
	return nil
}

func SetGrid(appId, spaceId, gridId string, g *Grid) {
	if app, ok := AppInfoL[appId]; ok {
		if space, ok := app.SpaceM[spaceId]; ok {
			space.GridM[gridId] = g
		} else {
			log.Errorln("not found space:", spaceId)
		}
	} else {
		log.Errorln("not found app:", appId)
	}
}

func RmGrid(appId, spaceId, gridId string) {
	if app, ok := AppInfoL[appId]; ok {
		if space, ok := app.SpaceM[spaceId]; ok {
			delete(space.GridM, gridId)
		}
	}
}
