package util

import (
	. "github.com/Unknwon/goconfig"
	log "github.com/Sirupsen/logrus"
)


var confFile = "conf.ini"
var Conf *Config = nil


const (
	SCT_DB = "db"
	SCT_HTTP = "http"
	SCT_TCP = "tcp"
)


type Config struct {
	c *ConfigFile
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
}
