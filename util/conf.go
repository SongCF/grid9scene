package util

import (
	. "github.com/Unknwon/goconfig"
	log "github.com/Sirupsen/logrus"
	"time"
	"strconv"
)


var confFile = "conf.ini"
var Conf *Config = nil

// will change by config file when start app
var (
	ReadBufSize = 2048
	WriteBufSize = 2048
	ReadDeadline time.Duration = 120 //second
)


const (
	SCT_DB = "db"
	SCT_HTTP = "http"
	SCT_TCP = "tcp"
)


type Config struct {
	c *ConfigFile
}


func (c *Config) GetInt(section, key string) int {
	val := c.Get(section, key)
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Errorf("config get int error: key=%s, val=%s, err=%v", key, val, err)
		return 0
	}
	return i
}
func (c *Config) Get(section, key string) string {
	// GetValue
	value, err := c.c.GetValue(section, key)
	if err != nil {
		log.Errorf("Error when get config value, err = %v", err)
		return ""
	}
	return value
}
func (c *Config) Gets(section string, keys []string) []string {
	ret := make([]string, len(keys))
	for i,k := range keys {
		ret[i] = c.Get(section, k)
	}
	return ret
}


func InitConf() {
	c, err := LoadConfigFile(confFile)
	CheckError(err)
	Conf = &Config{c:c}

	//
	loadTcpConfig()
}

func loadTcpConfig() {
	//read write buf
	ReadBufSize = Conf.GetInt(SCT_TCP, "read_buf")
	WriteBufSize = Conf.GetInt(SCT_TCP, "write_buf")
	//deadline
	ReadDeadline = time.Duration(Conf.GetInt(SCT_TCP, "deadline_time"))
}
