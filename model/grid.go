package model



type Grid struct {
	GridId string      // "x,y"
	X, Y   int         // x, y
	MsgBox chan []byte // message box
}

