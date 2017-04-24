package model

import (
	"database/sql"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	. "jhqc.com/songcf/scene/util"
	"strings"
)

var DB *sql.DB

func InitDB() {
	auth, err := Conf.Get(SCT_DB, "db_auth")
	CheckError(err)
	srv, err := Conf.Get(SCT_DB, "db_server")
	CheckError(err)
	db, err := Conf.Get(SCT_DB, "db_database")
	CheckError(err)
	//root:123456@tcp(139.198.5.219:3308)/db_scene_go?charset=utf8
	dst := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8", auth, srv, db)
	log.Println("init db: ", dst)
	mysql, err := sql.Open("mysql", dst)
	CheckError(err)
	DB = mysql
	maxOpen, err := Conf.GetInt(SCT_DB, "db_max_open_conn")
	CheckError(err)
	maxIdle, err := Conf.GetInt(SCT_DB, "db_max_idle_conn")
	CheckError(err)
	DB.SetMaxOpenConns(maxOpen)
	DB.SetMaxIdleConns(maxIdle)
	err = DB.Ping()
	CheckError(err)
}

func IsDuplicate(err error) bool {
	if err == nil {
		return false
	}
	str := err.Error()
	return strings.HasPrefix(str, "Error 1062: Duplicate entry")
}
