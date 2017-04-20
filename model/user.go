package model



type UserInfo struct {
	SpaceId  string  `json:"space_id"`
	GridId   string  `json:"grid_id"`
	PosX     float32 `json:"pos_x"`
	PosY     float32 `json:"pos_y"`
	Angle    float32 `json:"angle"`
	ExData   []byte  `json:"-"`
	MoveTime int32   `json:"-"`
}
