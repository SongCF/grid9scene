package util

import (
	"os"
	log "github.com/Sirupsen/logrus"
)


func CheckError(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
