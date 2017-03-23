package types

import (
	"net"
	"time"
)


type App struct {
	AppId string		        // "appId123456"
	SpaceList map[string]*Space      // spaceId : Space
	SessionList map[int32]*Session   // uid : Session
}


type Space struct {
	SpaceId string		     // "spaceId123456"
	GridWidth int32		     // 9-grid width
	GridHeight int32	     // 9-grid height
	GridList map[string]*Grid    // gridId : Grid
}


type Grid struct {
	GridId string     	     // "x,y"
	X, Y int		     // x, y
	ChanMsg chan []byte   	     // message box
}



type Session struct {
	AppId	     string
	Uid          int32
	ChanOut      chan []byte

	IP           net.IP
	Conn         net.Conn
	PacketCount  int32         //对包进行计数
	ConnectTime  time.Time

	UserData     *UserData
}

type UserData struct {
	SpaceId	     string
	GridId	     string
	PosX	     float32
	PosY	     float32
	Angle 	     float32
	ExData	     []byte
}

