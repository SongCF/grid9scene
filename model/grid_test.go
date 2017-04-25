package model

import "testing"

func TestCalcGridId(t *testing.T) {
	ret := CalcGridId(100.98, 23.5, 100, 100)
	if ret != "1,0" {
		t.Error("1.calc grid id error")
	}
	ret = CalcGridId(100.98, -23.5, 10, 10)
	if ret != "10,-3" {
		t.Errorf("2.calc grid id error, ret:%v", ret)
	}
}

func TestGetGridId(t *testing.T) {
	if GetGridId(1, -1) != "1,-1" {
		t.Error("get grid id error")
	}
}

func TestGetGridXY(t *testing.T) {
	x, y, err := GetGridXY("1,-1")
	if err != nil || x != 1 || y != -1 {
		t.Error("get grid xy error")
	}
}

func TestRoundGridAndSelf(t *testing.T) {
	ret := RoundGridAndSelf("3,3")
	if len(*ret) != 9 {
		t.Error("1.get round grid error")
	}
	checkHas(ret, []string{"2,2", "2,3", "2,4", "3,2", "3,3", "3,4", "4,2", "4,3", "4,4"}, t)

	ret = RoundGridAndSelf("0,0")
	if len(*ret) != 9 {
		t.Error("1.get round grid error")
	}
	checkHas(ret, []string{"-1,-1", "-1,0", "-1,1", "0,-1", "0,0", "0,1", "1,-1", "1,0", "1,1"}, t)

	//error
	ret = RoundGridAndSelf("0..0")
	if len(*ret) != 0 {
		t.Error("3.get round grid error")
	}
}

func checkHas(origin *[]string, dst []string, t *testing.T) {
	for _, g := range dst {
		has := false
		for _, o := range *origin {
			if g == o {
				has = true
				break
			}
		}
		if !has {
			t.Error("chack has false")
		}
	}
}
