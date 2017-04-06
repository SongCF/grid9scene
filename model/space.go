package model


type Space struct {
	SpaceId    string           // "spaceId123456"
	GridWidth  float32            // 9-grid width
	GridHeight float32            // 9-grid height
	GridM      map[string]*Grid // gridId : Grid
	MsgBox chan *InnerMsg `json:"-"`// message box
	Die chan struct{} `json:"-"`
}

func (s *Space) PostMsg(m *InnerMsg) {
	if s != nil && m != nil {
		select {
		case <- s.Die:
		case s.MsgBox <- m:
		}
	}
}

func (s *Space) Close() {
	if s != nil {
		select {
		case <- s.Die:
		default:
			close(s.Die)
		}
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

