package main

import (
	"time"
	"net"
)

type Session struct {
	ip           net.IP
	conn         net.Conn

	app_id       string
	user_id      int32

	in           chan []byte
	out          chan []byte
	ntf          chan []byte
	die          chan struct{} //会话关闭信号

	packet_count int32         //对包进行计数
	connect_time time.Time
}



func handleClient(conn net.Conn) {

}

func agent() {

}

func sender() {

}