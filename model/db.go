package model

import (
	log "github.com/Sirupsen/logrus"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"jhqc.com/songcf/scene/util"
	"strings"
)

var DB *sql.DB


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
	//rows, err := DB.Query("SELECT * FROM app;")
	//defer rows.Close()
	//util.CheckError(err)
	//for rows.Next() {
	//	var appid, name, key string
	//	if err := rows.Scan(&appid, &name, &key); err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println("appid:", appid, "\nname:", name, "\nkey:", key)
	//}
	//
	//sql := fmt.Sprintf("INSERT INTO app(app_id,name,private_key) VALUES('%s','%s','%s');", "a", "a", "a")
	//ret, err := DB.Exec(sql)
	//log.Println("ret:", ret)
	//log.Println("err:", err)

	//raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;", "1", "1")
	//var w, h float32
	//err = raw.Scan(&w, &h) // if empty, err = sql.ErrNoRows
	//if err != nil {
	//	log.Errorf("select grid w h error(%v:%v) = %v\n", 1, 1, err)
	//}
	//log.Infof("w:%v, h:%v", w, h)
}

func IsDuplicate(err error) bool {
	str := err.Error()
	return strings.HasPrefix(str, "Error 1062: Duplicate entry")
}