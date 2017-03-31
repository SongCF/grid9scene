package main

import (
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/global"
	. "jhqc.com/songcf/scene/gateway"
	. "jhqc.com/songcf/scene/model"
	"jhqc.com/songcf/scene/util"
	. "jhqc.com/songcf/scene/controller"
)



func main() {
	log.SetLevel(log.DebugLevel)

	go HandleSignal()

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
	util.CheckError(err)
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