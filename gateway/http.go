package gateway

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/SongCF/scene/model"
	"github.com/SongCF/scene/pb"
	. "github.com/SongCF/scene/util"
	"net/http"
	"strconv"
)

func HttpServer() {
	addr, err := Conf.Get(SCT_HTTP, "http_server")
	CheckError(err)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/app/{aid}", handleCreateApp).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}", handleDeleteApp).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", handleCreateSpace).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", handleDeleteSpace).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/user/{uid:[0-9]+}/pos", handleQueryPos).Methods("GET")
	//r.StrictSlash(true)
	log.Println("http server listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

//"http://127.0.0.1:9911/api/v1/app/1?user_id=1&token=abc123"
func handleCreateApp(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	appId := vars["aid"]
	eInfo := CreateApp(appId, "test", "test")
	rsp := pb.ErrSuccess
	if eInfo != nil {
		rsp = eInfo
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("json encode failed, rsp=%v", rsp)
	}
	w.Write(b)
}

func handleDeleteApp(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	appId := vars["aid"]
	eInfo := DeleteApp(appId)
	rsp := pb.ErrSuccess
	if eInfo != nil {
		rsp = eInfo
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("json encode failed, rsp=%v", rsp)
	}
	w.Write(b)
}

//"http://127.0.0.1:9911/api/v1/app/1/space/1?user_id=1&token=abc123&grid_width=10&grid_height=10"
func handleCreateSpace(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	appId := vars["aid"]
	spaceId := vars["sid"]
	r.ParseForm()
	width, err1 := strconv.ParseFloat(r.Form.Get("grid_width"), 32)
	height, err2 := strconv.ParseFloat(r.Form.Get("grid_height"), 32)

	rsp := pb.ErrSuccess
	if err1 != nil || err2 != nil {
		rsp = pb.ErrMsgFormat
	} else {
		eInfo := CreateSpace(appId, spaceId, float32(width), float32(height))
		if eInfo != nil {
			rsp = eInfo
		}
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("json encode failed, rsp=%v", rsp)
	}
	w.Write(b)
}

func handleDeleteSpace(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	appId := vars["aid"]
	spaceId := vars["sid"]
	eInfo := DeleteSpace(appId, spaceId)
	rsp := pb.ErrSuccess
	if eInfo != nil {
		rsp = eInfo
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("json encode failed, rsp=%v", rsp)
	}
	w.Write(b)
}

//"http://127.0.0.1:9911/api/v1/app/1/user/1/pos?user_id=1&token=abc123"
func handleQueryPos(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	appId := vars["aid"]
	queryUidStr := vars["uid"]

	rsp := pb.ErrSuccess
	queryUid, err1 := strconv.Atoi(queryUidStr)
	if err1 != nil {
		rsp = pb.ErrMsgFormat
	} else {
		conn, err := CCPool.Get()
		if err != nil {
			log.Errorf("handleQueryPos(user[%v:%v]) CCPool:Get error(%v)", appId, queryUid, err)
			rsp = pb.ErrServerBusy
		} else {
			defer CCPool.Put(conn)
			userInfo, e := GetUserInfo(appId, int32(queryUid), conn)
			if e != nil {
				rsp = e
			} else {
				rsp.Ex = fmt.Sprintf("{\"space_id\":\"%v\", \"x\":%v, \"y\":%v, \"angle\":%v}",
					userInfo.SpaceId, userInfo.PosX, userInfo.PosY, userInfo.Angle)
			}
		}
	}
	b, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("json encode failed, rsp=%v", rsp)
	}
	w.Write(b)
}
