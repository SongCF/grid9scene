package model

import (
	"fmt"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestDB(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	var app_id = "test_app_id"
	var space_id = "test_space_id"

	CreateApp(app_id, "test_name", "test_key")

	rows, err := DB.Query("SELECT * FROM app;")
	defer rows.Close()
	if err != nil {
		t.Errorf("select app error:%v", err)
	}
	for rows.Next() {
		var appid, name, key string
		if err := rows.Scan(&appid, &name, &key); err != nil {
			t.Errorf("scan app error:%v", err)
		}
	}

	sql := fmt.Sprintf("INSERT INTO app(app_id,name,private_key) VALUES('%s','%s','%s');", "a", "a", "a")
	_, err = DB.Exec(sql)
	if !IsDuplicate(err) && err != nil {
		t.Errorf("insert err:%v", err)
	}

	CreateSpace(app_id, space_id, float32(10), float32(10))

	raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;", app_id, space_id)
	var w, h float32
	err = raw.Scan(&w, &h) // if empty, err = sql.ErrNoRows
	if err != nil {
		t.Errorf("select grid err:%v", err)
	}
}
