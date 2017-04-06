package model


type App struct {
	SpaceM map[string]*Space      // spaceId : Space
	SessionM map[int32]*Session // uid : Session
	MsgBox chan *InnerMsg `json:"-"`// message box
	Die chan struct{} `json:"-"`
}

func (app *App) PostMsg(msg *InnerMsg) {
	if app != nil  && msg != nil {
		select {
		case <- app.Die:
		case app.MsgBox <- msg:
		}
	}
}

func (app *App) Close() {
	if app != nil {
		select {
		case <- app.Die:
		default:
			close(app.Die)
		}
	}
}


func HasApp(appId string) bool {
	if _,ok := AppL[appId]; ok {
		return true
	}
	return false
}


func GetApp(appId string) *App {
	if app, ok := AppL[appId]; ok {
		return app
	}
	return nil
}
