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
	vL := Conf.Gets(SCT_DB, []string{"db_server", "db_auth", "db_database"})
	//root:123456@tcp(139.198.5.219:3308)/db_scene_go?charset=utf8
	dst := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8", vL[1], vL[0], vL[2])
	log.Println("init db: ", dst)
	mysql, err := sql.Open("mysql", dst)
	CheckError(err)
	DB = mysql
	DB.SetMaxOpenConns(Conf.GetInt(SCT_DB, "db_max_open_conn", 2000))
	DB.SetMaxIdleConns(Conf.GetInt(SCT_DB, "db_max_idle_conn", 1000))
	err = DB.Ping()
	CheckError(err)

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
