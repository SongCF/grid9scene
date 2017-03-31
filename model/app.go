package model


type App struct {
	SpaceM map[string]*Space      // spaceId : Space
	SessionM map[int32]*Session // uid : Session
	MsgBox chan *InnerMsg `json:"-"`// message box
}

func (app *App) PostMsg(msg *InnerMsg) {
	if app != nil && app.MsgBox != nil && msg != nil {
		app.MsgBox <- msg
	}

}

func (app *App) Close() {
	if app != nil && app.MsgBox != nil {
		close(app.MsgBox)
		app.MsgBox = nil
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
