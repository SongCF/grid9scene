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

var (
	ZkConn     *zk.Conn
	zkTcpKey   string
	zkTcpData  []byte
	zkHttpKey  string
	zkHttpData []byte
)

// Register to zk
func Register() {
	ZkConn = connect()

	log.Info("Registering to Zookeeper")
	regTcpUri, err := Conf.Get(SCT_ZK, "zk_reg_tcp_uri")
	CheckError(err)
	regHttpUri, err := Conf.Get(SCT_ZK, "zk_reg_http_uri")
	CheckError(err)
	regTcpAddr, err := Conf.Get(SCT_ZK, "zk_reg_tcp_addr")
	CheckError(err)
	regHttpAddr, err := Conf.Get(SCT_ZK, "zk_reg_http_addr")
	CheckError(err)
	if regTcpUri == "" || regHttpUri == "" || regTcpAddr == "" || regHttpAddr == "" {
		panic("register zookeeper param error")
	}

	//tcp
	zkTcpKey = fmt.Sprintf("%v/%v", regTcpUri, regTcpAddr)
	zkTcpData, err = json.Marshal(zkData{URI: fmt.Sprintf("http://%v", regTcpAddr)})
	CheckError(err)
	createNode(zkTcpKey, zkTcpData)
	//http
	zkHttpKey = fmt.Sprintf("%v/%v", regHttpUri, regHttpAddr)
	zkHttpData, err = json.Marshal(zkData{URI: fmt.Sprintf("http://%v", regHttpAddr)})
	CheckError(err)
	createNode(zkHttpKey, zkHttpData)

	log.Info("Register To Zookeeper OK")
}

func Unregister() {
	ZkConn.Delete(zkTcpKey, -1)
	ZkConn.Delete(zkHttpKey, -1)
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
		_, err := ZkConn.Create(node, []byte(node), int32(0), acl)
		if err != nil && err != zk.ErrNodeExists {
			log.Errorf("create node:%v, err:%v", node, err)
			panic(err)
		}
	}
	if err := ZkConn.Delete(regKey, -1); err != nil {
		log.Info("Skip deleting zk registration ", regKey)
	} else {
		log.Info("zk registration ", regKey, " deleted")
	}
	_, err := ZkConn.Create(regKey, data, zk.FlagEphemeral, acl)
	CheckError(err)
}

// GetServices
// @param name:   tcp / http
func GetServices(name string) [][]byte {
	if name == "" {
		log.Info("ZK GetServices name is none")
		return [][]byte{}
	}
	key := ""
	var err error
	if name == "tcp" {
		key, err = Conf.Get(SCT_ZK, "zk_reg_tcp_uri")
		CheckError(err)
	} else if name == "http" {
		key, err = Conf.Get(SCT_ZK, "zk_reg_http_uri")
		CheckError(err)
	} else {
		log.Infof("ZK GetServices name(%v) not support", name)
		return [][]byte{}
	}
	children, _, _ := ZkConn.Children(key)

	retL := [][]byte{}
	for _, c := range children {
		data, _, _ := ZkConn.Get(fmt.Sprintf("%v/%v", key, c))
		retL = append(retL, data)
		log.Info(key+"/"+c, ", Data:", string(data))
	}
	return retL
}

func connect() *zk.Conn {
	serverStr, err := Conf.Get(SCT_ZK, "zk_servers")
	CheckError(err)
	auth, err := Conf.Get(SCT_ZK, "zk_auth")
	CheckError(err)
	zkTimeout, err := Conf.GetInt(SCT_ZK, "zk_timeout")
	CheckError(err)
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
