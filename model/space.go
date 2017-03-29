package model


type Space struct {
	SpaceId    string           // "spaceId123456"
	GridWidth  float32            // 9-grid width
	GridHeight float32            // 9-grid height
	GridM      map[string]*Grid // gridId : Grid
	MsgBox chan *InnerMsg // message box
}

func (s *Space) PostMsg(m *InnerMsg) {
	if s != nil && s.MsgBox != nil && m != nil {
		s.MsgBox <- m
	}
}

func (s *Space) Close() {
	if s != nil && s.MsgBox != nil {
		close(s.MsgBox)
		s.MsgBox = nil
	}
}


func GetSpace(appId, spaceId string) *Space {
	if app, ok := AppL[appId]; ok {
		if s, ok := app.SpaceM[spaceId]; ok {
			return s
		}
	}
	return nil
}