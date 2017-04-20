package util

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/Unknwon/goconfig"
	"time"
)

var confFile = "conf.ini"
var Conf *Config = nil

// will change by config file when start app
var (
	ReadBufSize                = 2048
	WriteBufSize               = 2048
	ReadDeadline time.Duration = 120 //second
)

const (
	SCT_DB   = "db"
	SCT_CACHE   = "cache"
	SCT_HTTP = "http"
	SCT_TCP  = "tcp"
	SCT_ZK = "zookeeper"
)

type Config struct {
	c *ConfigFile
}

func (c *Config) GetInt(section, key string, dft int) int {
	val, err := c.c.Int(section, key)
	if err != nil {
		log.Errorf("config get int error: key=[%v]%s, val=%s, err=%v", section, key, val, err)
		return dft
	}
	return val
}
func (c *Config) Get(section, key string) string {
	// GetValue
	value, err := c.c.GetValue(section, key)
	if err != nil {
		log.Errorf("Error when get config value([%v]%v), err = %v", section, key, err)
		return ""
	}
	return value
}
func (c *Config) Gets(section string, keys []string) []string {
	ret := make([]string, len(keys))
	for i, k := range keys {
		ret[i] = c.Get(section, k)
	}
	return ret
}

func InitConf() {
	c, err := LoadConfigFile(confFile)
	CheckError(err)
	Conf = &Config{c: c}

	//
	loadTcpConfig()
}

func loadTcpConfig() {
	//read write buf
	ReadBufSize = Conf.GetInt(SCT_TCP, "tcp_read_buf", 2048)
	WriteBufSize = Conf.GetInt(SCT_TCP, "tcp_write_buf", 2048)
	//deadline
	ReadDeadline = time.Duration(Conf.GetInt(SCT_TCP, "tcp_dead_time", 120)) * time.Second
}
