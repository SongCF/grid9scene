package model

import (
	"fmt"
	"jhqc.com/songcf/scene/_test"
	. "jhqc.com/songcf/scene/util"
	"testing"
)

func TestDB(t *testing.T) {
	InitConfTest("../conf.ini")
	InitDB()

	CreateApp(_test.T_APP_ID, "test_name", "test_key")

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

	CreateSpace(_test.T_APP_ID, _test.T_SPACE_ID, float32(10), float32(10))

	raw := DB.QueryRow("SELECT grid_width,grid_height FROM space WHERE app_id=? and space_id=?;",
		_test.T_APP_ID, _test.T_SPACE_ID)
	var w, h float32
	err = raw.Scan(&w, &h) // if empty, err = sql.ErrNoRows
	if err != nil {
		t.Errorf("select grid err:%v", err)
	}
}
