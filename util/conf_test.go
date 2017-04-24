package util

import (
	. "github.com/Unknwon/goconfig"
	"testing"
)

func prepareConf() {
	//InitConf()
	c, err := LoadConfigFile("../" + confFile)
	CheckError(err)
	Conf = &Config{c: c}
}

func TestConfig_Get(t *testing.T) {
	prepareConf()
	_, err := Conf.Get(SCT_DB, "db_server")
	if err != nil {
		t.Errorf("Conf get db_server error:%v", err)
	}
	_, err = Conf.GetInt(SCT_DB, "db_max_open_conn")
	if err != nil {
		t.Errorf("Conf get db_max_open_conn error:%v", err)
	}
}
