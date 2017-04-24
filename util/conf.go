package util

import (
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
	SCT_DB    = "db"
	SCT_CACHE = "cache"
	SCT_HTTP  = "http"
	SCT_TCP   = "tcp"
	SCT_ZK    = "zookeeper"
)

type Config struct {
	c *ConfigFile
}

func (c *Config) GetInt(section, key string) (int, error) {
	return c.c.Int(section, key)
}
func (c *Config) Get(section, key string) (string, error) {
	// GetValue
	return c.c.GetValue(section, key)
}

func InitConf() {
	c, err := LoadConfigFile(confFile)
	CheckError(err)
	Conf = &Config{c: c}

	//
	loadTcpConfig()
}

// use for test
func InitConfTest(file string) {
	c, err := LoadConfigFile(file)
	CheckError(err)
	Conf = &Config{c: c}
}

func loadTcpConfig() {
	var err error
	//read write buf
	ReadBufSize, err = Conf.GetInt(SCT_TCP, "tcp_read_buf")
	CheckError(err)
	WriteBufSize, err = Conf.GetInt(SCT_TCP, "tcp_write_buf")
	CheckError(err)
	//deadline
	dt, err := Conf.GetInt(SCT_TCP, "tcp_dead_time")
	CheckError(err)
	ReadDeadline = time.Duration(dt) * time.Second
}
