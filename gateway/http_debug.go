package gateway

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"github.com/gorilla/mux"
	. "jhqc.com/songcf/scene/model"
	"encoding/json"
	. "jhqc.com/songcf/scene/util"
	"fmt"
)

func StartPProf() {
	addr := Conf.Get(SCT_HTTP, "pprof")
	log.Println("pprof listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}


func StartStats() {
	addr := Conf.Get(SCT_HTTP, "stat")
	r := mux.NewRouter()
	r.HandleFunc("/stat/{type}", handleStat).Methods("GET")
	//r.StrictSlash(true)
	log.Println("stat listening on: ", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func handleStat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryType := vars["type"]

	res := []byte{}
	switch queryType {
	case "cache":
		s, err := json.MarshalIndent(AppL, "", " ")
		if err != nil {
			eStr := fmt.Sprintf("Marshal AppL to json error: %v", err)
			w.Write([]byte(eStr))
			return
		}
		res = s
	default:
		res = []byte("unknown stats type")
		log.Errorln("query stats error type: ", queryType)
	}
	w.Write(res)
}
