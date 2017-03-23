package util

import (
	"os"
	log "github.com/Sirupsen/logrus"
)


func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
