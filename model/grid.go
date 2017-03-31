package model

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)


type Grid struct {
	GridId string      // "x,y"
	UidM map[int32]bool  // grid uid list
	MsgBox chan interface{} `json:"-"`// message box
}

func (g *Grid) PostMsg(m interface{}) {
	if g != nil && g.MsgBox != nil && m != nil {
		g.MsgBox <- m
	}
}
func (g *Grid) Close() {
	if g != nil && g.MsgBox != nil {
		close(g.MsgBox)
		g.MsgBox = nil
	}
}



func CalcGridId(posX, posY, w, h float32) string {
	x := int32(posX / w)
	y := int32(posY / h)
	return GetGridId(x, y)
}

func GetGridId(x, y int32) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func GetGridXY(gridId string) (int32, int32, error) {
	var x, y int32
	_, err := fmt.Sscanf(gridId, "%d,%d", x, y)
	return x, y, err
}

//include self
func RoundGridAndSelf(gridId string) *[]string {
	gridIdL := []string{}
	x, y, err := GetGridXY(gridId)
	if err != nil {
		log.Errorf("get gridxy failed, gridid=%v", gridId)
		return &gridIdL
	}
	gridIdL = append(gridIdL, GetGridId(x-1, y-1))
	gridIdL = append(gridIdL, GetGridId(x-1, y))
	gridIdL = append(gridIdL, GetGridId(x-1, y+1))
	gridIdL = append(gridIdL, GetGridId(x, y-1))
	gridIdL = append(gridIdL, GetGridId(x, y))
	gridIdL = append(gridIdL, GetGridId(x, y+1))
	gridIdL = append(gridIdL, GetGridId(x+1, y-1))
	gridIdL = append(gridIdL, GetGridId(x+1, y))
	gridIdL = append(gridIdL, GetGridId(x+1, y+1))
	return &gridIdL
}

func SubGrids(arr1, arr2 *[]string) *[]string {
	arr := []string{}
	for i := 0; i < len(*arr1); i++ {
		b := true
		for j := 0; j < len(*arr2); j++ {
			if (*arr1)[i] == (*arr2)[j] {
				b = false
				break
			}
		}
		if b {
			arr = append(arr, (*arr1)[i])
		}
	}
	return &arr
}
