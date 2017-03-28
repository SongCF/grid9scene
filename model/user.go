package model


type UserInfo struct {
	SpaceId	     string
	GridId	     string
	PosX	     float32
	PosY	     float32
	Angle 	     float32
	ExData	     []byte
	MoveTime     int32
}

