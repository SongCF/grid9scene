package main

import (
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
)

func httpServer() {
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
	createApp("1", "1", "1")
	log.Info("test!")
	w.Write([]byte("createApp\n"))
}

func handleDeleteApp(w http.ResponseWriter, r *http.Request) {
	deleteApp("1")
	w.Write([]byte("deleteApp\n"))
}

func handleCreateSpace(w http.ResponseWriter, r *http.Request) {
	createSpace("1", "1", 1, 1)
	w.Write([]byte("createSpace\n"))
}

func handleDeleteSpace(w http.ResponseWriter, r *http.Request) {
	deleteSpace("1", "1")
	w.Write([]byte("deleteSpace\n"))
}

func handleQueryPos(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("queryPos\n"))
}