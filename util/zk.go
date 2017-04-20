package util

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

type zkData struct {
	URI string `json:"uri"`
}

var conn *zk.Conn

// Register to zk
func Register() {
	conn = connect()

	log.Info("Registering to Zookeeper")
	regTcpUri := Conf.Get(SCT_ZK, "zk_reg_tcp_uri")
	regHttpUri := Conf.Get(SCT_ZK, "zk_reg_http_uri")
	regTcpAddr := Conf.Get(SCT_ZK, "zk_reg_tcp_addr")
	regHttpAddr := Conf.Get(SCT_ZK, "zk_reg_http_addr")
	if regTcpUri == "" || regHttpUri == "" || regTcpAddr == "" || regHttpAddr == "" {
		panic("register zookeeper param error")
	}

	//tcp
	zkTcpKey := fmt.Sprintf("%v/%v", regTcpUri, regTcpAddr)
	zkTcpData, err := json.Marshal(zkData{URI: fmt.Sprintf("http://%v", regTcpAddr)})
	CheckError(err)
	createNode(zkTcpKey, zkTcpData)
	//http
	zkHttpKey := fmt.Sprintf("%v/%v", regHttpUri, regHttpAddr)
	zkHttpData, err := json.Marshal(zkData{URI: fmt.Sprintf("http://%v", regHttpAddr)})
	CheckError(err)
	createNode(zkHttpKey, zkHttpData)

	log.Info("Register To Zookeeper OK")
}

func createNode(regKey string, data []byte) {
	topList := strings.Split(regKey, "/")
	parentList := []string{}
	tag := ""
	for _, p := range topList {
		if p != "" {
			tag = tag + "/" + p
			parentList = append(parentList, tag)
		}
	}
	acl := zk.WorldACL(zk.PermAll)
	for _, node := range parentList {
		if node == regKey {
			continue
		}
		// Create as a persistent znode, since the zk package lacks the constant for the
		// persistent type znode, thus the magic number 0 here.
		_, err := conn.Create(node, []byte(node), int32(0), acl)
		if err != nil && err != zk.ErrNodeExists {
			log.Errorf("create node:%v, err:%v", node, err)
			panic(err)
		}
	}
	if err := conn.Delete(regKey, -1); err != nil {
		log.Info("Skip deleting zk registration ", regKey)
	} else {
		log.Info("zk registration ", regKey, " deleted")
	}
	_, err := conn.Create(regKey, data, zk.FlagEphemeral, acl)
	CheckError(err)
}

// GetServices
// @param name:   tcp / http
func GetServices(name string) {
	if name == "" {
		log.Info("ZK GetServices name is none")
		return
	}
	key := ""
	if name == "tcp" {
		key = Conf.Get(SCT_ZK, "zk_reg_tcp_uri")
	} else if name == "http" {
		key = Conf.Get(SCT_ZK, "zk_reg_http_uri")
	} else {
		log.Infof("ZK GetServices name(%v) not support", name)
		return
	}
	children, _, _ := conn.Children(key)

	for _, c := range children {
		data, _, _ := conn.Get(fmt.Sprintf("%v/%v", key, c))
		log.Info(key+"/"+c, ", Data:", string(data))
	}
}

func connect() *zk.Conn {
	serverStr := Conf.Get(SCT_ZK, "zk_servers")
	auth := Conf.Get(SCT_ZK, "zk_auth")
	zkTimeout := Conf.GetInt(SCT_ZK, "zk_timeout", 10)
	if serverStr == "" || auth == "" {
		panic("None zookeeper addr or auth")
	}
	zkServers := strings.Split(serverStr, ",")

	log.Infoln("Connecting Zookeeper ...", zkServers)
	conn, _, err := zk.Connect(zkServers, time.Duration(zkTimeout)*time.Second)
	if err != nil {
		log.Errorln("Establish connect to Zookeeper error: ", err)
		panic(err)
	}

	log.Infoln("zk auth: ", auth)
	err = conn.AddAuth("digest", []byte(auth))
	if err != nil {
		log.Errorf("Zookeeper Auth failed, auth=%v, err=%v", auth, err)
		panic(err)
	}
	return conn
}
