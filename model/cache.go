package model


import (
	log "github.com/Sirupsen/logrus"
)




// mem cache

var (
	AppL = make(map[string]*App)
)



func GetGrid(appId, spaceId, gridId string) *Grid {
	if app, ok := AppL[appId]; ok {
		if space, ok := app.SpaceM[spaceId]; ok {
			if grid, ok := space.GridM[gridId]; ok {
				return grid
			}
		}
	}
	return nil
}

func SetGrid(appId, spaceId, gridId string, g *Grid) {
	if app, ok := AppL[appId]; ok {
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
	if app, ok := AppL[appId]; ok {
		if space, ok := app.SpaceM[spaceId]; ok {
			delete(space.GridM, gridId)
		}
	}
}
