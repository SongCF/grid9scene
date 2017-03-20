package main

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
)

func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/app/{aid}", createApp).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}", deleteApp).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", createSpace).Methods("POST")
	r.HandleFunc("/api/v1/app/{aid}/space/{sid}", deleteSpace).Methods("DELETE")
	r.HandleFunc("/api/v1/app/{aid}/user/{uid:[0-9]+}/pos", queryPos).Methods("GET")
	//r.StrictSlash(true)
	log.Fatal(http.ListenAndServe(":9911", r))
}


func createApp(w http.ResponseWriter, r *http.Request) {
	log.Info("test!")
	w.Write([]byte("createApp\n"))
}

func deleteApp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteApp\n"))
}

func createSpace(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createSpace\n"))
}

func deleteSpace(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteSpace\n"))
}

func queryPos(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("queryPos\n"))
}