package model

type InnerMsg struct {
	Id      int32
	Cb      chan struct{}
	AppId   string
	SpaceId string
	GridId  string
}

// inner msg
const (
	//IMSG_START_APP int32 = 1
	IMSG_START_SPACE int32 = 2
	IMSG_START_GRID  int32 = 3
)
