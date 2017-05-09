package _test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jhqc.com/songcf/scene/pb"
	"net/http"
	"strings"
	"testing"
)

func TestAllAPI(addr string, t *testing.T) {
	fmt.Println("http addr: ", addr)
	var err error
	//清理
	err = deleteApp(addr)
	if err != nil {
		t.Fatal(err)
	}

	//无应用时创建场景会失败
	err = createSpace(addr)
	if err.Error() != pb.ErrAppNotExist.Desc {
		t.Fatal(err)
	}

	//创建应用成功
	err = createApp(addr)
	if err != nil {
		t.Fatal(err)
	}
	//已存在
	err = createApp(addr)
	if err.Error() != pb.ErrAppAlreadyExist.Desc {
		t.Fatal(err)
	}

	//没有时 删除是成功
	err = deleteSpace(addr)
	if err != nil {
		t.Fatal(err)
	}
	//创建场景成功
	err = createSpace(addr)
	if err != nil {
		t.Fatal(err)
	}
	//已存在
	err = createSpace(addr)
	if err.Error() != pb.ErrSpaceAlreadyExist.Desc {
		t.Fatal(err)
	}

	//查询位置(用户不在线)
	err = queryPos(addr)
	if err.Error() != pb.ErrUserOffline.Desc {
		t.Fatal(err)
	}

	//删除场景
	err = deleteSpace(addr)
	if err != nil {
		t.Fatal(err)
	}

	//删除应用
	err = deleteApp(addr)
	if err != nil {
		t.Fatal(err)
	}

	//查询位置（应用不存在）
	err = queryPos(addr)
	if err.Error() != pb.ErrAppNotExist.Desc {
		t.Fatal(err)
	}
}

func initAppSpace(addr string) {
	err := createApp(addr)
	if err != nil && err.Error() != pb.ErrAppAlreadyExist.Desc {
		panic(err)
	}
	err = createSpace(addr)
	if err != nil && err.Error() != pb.ErrSpaceAlreadyExist.Desc {
		panic(err)
	}
}

func createApp(addr string) error {
	fmt.Println("createApp...")
	body := strings.NewReader("")
	uri := fmt.Sprintf("%s/api/v1/app/%s", addr, T_APP_ID)
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Printf("post create app, req error:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return parseResp(resp.Body)
}

func deleteApp(addr string) error {
	fmt.Println("deleteApp...")
	client := &http.Client{}
	body := strings.NewReader("")
	uri := fmt.Sprintf("%s/api/v1/app/%s", addr, T_APP_ID)
	req, err := http.NewRequest("DELETE", uri, body)
	if err != nil {
		fmt.Printf("(delete app) new request error:%v\n", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("(delete app) client do req error:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return parseResp(resp.Body)
}

func createSpace(addr string) error {
	fmt.Println("createSpace...")
	body := strings.NewReader("")
	uri := fmt.Sprintf("%s/api/v1/app/%s/space/%s?grid_width=%v&grid_height=%v",
		addr, T_APP_ID, T_SPACE_ID, T_GRID_W, T_GRID_H)
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Printf("post create space, req error:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return parseResp(resp.Body)
}

func deleteSpace(addr string) error {
	fmt.Println("deleteSpace...")
	client := &http.Client{}
	body := strings.NewReader("")
	uri := fmt.Sprintf("%s/api/v1/app/%s/space/%s", addr, T_APP_ID, T_SPACE_ID)
	req, err := http.NewRequest("DELETE", uri, body)
	if err != nil {
		fmt.Printf("(delete space) new request error:%v\n", err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("(delete space) client do req error:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return parseResp(resp.Body)
}

func queryPos(addr string) error {
	fmt.Println("queryPos...")
	uri := fmt.Sprintf("%s/api/v1/app/%s/user/%d/pos", addr, T_APP_ID, T_USER_ID)
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Printf("get user pos, req error:%v\n", err)
		return err
	}
	defer resp.Body.Close()
	return parseResp(resp.Body)
}

func parseResp(body io.Reader) error {
	rBody, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Printf("read body error:%v\n", err)
		return err
	}
	r := &pb.ErrInfo{}
	err = json.Unmarshal(rBody, r)
	if err != nil {
		fmt.Printf("parse body error:%v\n", err)
		return err
	}
	fmt.Printf("parse resp: %v\n", r)
	if r.Id != pb.ErrSuccess.Id {
		return errors.New(r.Desc)
	}
	return nil
}
