package gateway

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/SongCF/scene/model"
	. "github.com/SongCF/scene/util"
	"net/http"
	_ "net/http/pprof"
	"github.com/SongCF/scene/_test"
	"github.com/SongCF/scene/rpc"
	"github.com/SongCF/scene/pb"
	"github.com/golang/protobuf/proto"
)

func StartPProf() {
	addr, err := Conf.Get(SCT_HTTP, "pprof_server")
	CheckError(err)
	log.Println("pprof listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func StartStats() {
	addr, err := Conf.Get(SCT_HTTP, "stat_server")
	CheckError(err)
	r := mux.NewRouter()
	r.HandleFunc("/stat/{type}", handleStat).Methods("GET")
	//r.StrictSlash(true)
	log.Println("stat listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

//http://127.0.0.1:9912/stat/session
func handleStat(w http.ResponseWriter, r *http.Request) {
	defer RecoverPanic()
	vars := mux.Vars(r)
	queryType := vars["type"]

	res := []byte{}
	switch queryType {
	case "session":
		m := GetAllSession()
		s, err := json.MarshalIndent(m, "", " ")
		if err != nil {
			eStr := fmt.Sprintf("Marshal AppL to json error: %v", err)
			w.Write([]byte(eStr))
			return
		}
		res = s
	case "avg_msg":
		ssNum, pkCount, timeAll := GetAllSessionMsgAvg()
		avgPk := 0
		avgTime := 0
		if ssNum > 0 {
			avgPk = pkCount / ssNum
			avgTime = timeAll / ssNum
		}
		res = []byte(fmt.Sprintf("session num: %v,\npacket_count: %v,\ntime_all: %v,\navg_pk: %v,\navg_time: %v\n",
			ssNum, pkCount, timeAll, avgPk, avgTime))
	case "zk":
		zkTcp := GetServices("tcp")
		zkHttp := GetServices("http")
		tcpStr := ""
		httpStr := ""
		for _, t := range zkTcp {
			tcpStr += string(t) + "\n"
		}
		for _, h := range zkHttp {
			httpStr += string(h) + "\n"
		}
		res = []byte(fmt.Sprintf("tcp:\n%s\nhttp:\n%s", tcpStr, httpStr))
	case "test":
		m, err := proto.Marshal(&pb.LeaveAck{})
		if err != nil {
			res = []byte("Error: Marshal payload failed!")
		} else {
			r := rpc.TestRPC(":9902", _test.T_APP_ID, _test.T_USER_ID, pb.CmdLeaveAck, m)
			s, err := json.MarshalIndent(r, "", " ")
			if err != nil {
				eStr := fmt.Sprintf("Marshal AppL to json error: %v", err)
				res = []byte(eStr)
			} else {
				res = s
			}
		}
	default:
		res = []byte("unknown stats type")
		log.Errorln("query stats error type: ", queryType)
	}
	w.Write(res)
}
