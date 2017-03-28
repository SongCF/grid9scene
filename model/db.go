package model

import (
	log "github.com/Sirupsen/logrus"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"jhqc.com/songcf/scene/util"
)

var DB *sql.DB


// mysql table
const (
	TBL_APP = "app"
	TBL_SPACE = "space"
	TBL_LAST_SPACE = "last_space"
	TBL_LAST_POS = "last_pos"
)

func InitDB() {
	log.Info("init db")
	mysql, err := sql.Open("mysql", "root:123456@tcp(192.168.31.216:3306)/scene_db?charset=utf8")
	if err != nil {
		panic(err)
	}
	DB = mysql
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	err = DB.Ping()
	util.CheckError(err)

	////test
	//rows, err := db.Query("SELECT * FROM app;")
	//defer rows.Close()
	//checkError(err)
	//for rows.Next() {
	//	var appid, name, key string
	//	if err := rows.Scan(&appid, &name, &key); err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println("appid:", appid, "\nname:", name, "\nkey:", key)
	//}
}