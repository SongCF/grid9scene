package util

import (
	log "github.com/Sirupsen/logrus"
)

func CheckError(err error) {
	if err != nil {
		log.Errorf("CheckError panic: %v", err)
		panic(err)
	}
}
