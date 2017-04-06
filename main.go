package main

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/controller"
	. "jhqc.com/songcf/scene/gateway"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/model"
	. "jhqc.com/songcf/scene/util"
	"os"
)

func main() {
	defer RecoverPanic()

	log.SetLevel(log.DebugLevel)

	//
	wd, err := os.Getwd()
	CheckError(err)
	log.Println("work dir: ", wd)

	go HandleSignal()

	//load config
	InitConf()

	InitDB()
	loadAppTbl()

	go HttpServer()

	go TcpServer()

	go StartPProf()
	go StartStats()

	initZK()

	select {}
}

func loadAppTbl() {
	rows, err := DB.Query("SELECT app_id FROM app;")
	defer rows.Close()
	CheckError(err)
	for rows.Next() {
		var appId string
		if err := rows.Scan(&appId); err != nil {
			log.Fatal(err)
		}
		StartApp(appId)
	}
}

func initZK() {
	// TODO zookeeper
}
