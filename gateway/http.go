package gateway

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
	. "jhqc.com/songcf/scene/controller"
)

func HttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/app/{aid}", handleCreateApp).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}", handleDeleteApp).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", handleCreateSpace).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", handleDeleteSpace).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/user/{uid:[0-9]+}/pos", handleQueryPos).Methods("GET")
	//r.StrictSlash(true)
	log.Fatal(http.ListenAndServe(":9911", r))
}


func handleCreateApp(w http.ResponseWriter, r *http.Request) {
	appId := "1"
	err := CreateApp(appId, "1", "1")
	if err != nil {
		//TODO
		w.Write([]byte("createApp failed\n"))
		return
	}
	w.Write([]byte("createApp success\n"))
}

func handleDeleteApp(w http.ResponseWriter, r *http.Request) {
	err := DeleteApp("1")
	if err != nil {
		//TODO
		w.Write([]byte("deleteApp failed\n"))
		return
	}
	w.Write([]byte("deleteApp\n"))
}

func handleCreateSpace(w http.ResponseWriter, r *http.Request) {
	err := CreateSpace("1", "1", 1, 1)
	if err != nil {
		//TODO
		w.Write([]byte("createSpace failed\n"))
		return
	}
	w.Write([]byte("createSpace\n"))
}

func handleDeleteSpace(w http.ResponseWriter, r *http.Request) {
	err := DeleteSpace("1", "1")
	if err != nil {
		//TODO
		w.Write([]byte("deleteSpace failed\n"))
		return
	}
	w.Write([]byte("deleteSpace\n"))
}

func handleQueryPos(w http.ResponseWriter, r *http.Request) {
	//TODO
	w.Write([]byte("queryPos\n"))
}