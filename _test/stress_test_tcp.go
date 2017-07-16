package _test

import (
	"fmt"
	"github.com/SongCF/scene/pb"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

//模拟客户端行为，以tcp链接服务器，每隔 min_interval ~ max_interval 秒像服务器发送一条请求，
//链接服务器后，login 然后 join T_SPACE_ID 场景中，之后只发送 move_req, broadcast_req 两种消息，
//每个 leave_time 秒，下线部分客户端，然后重新上线

type Client struct {
	R chan []byte
	W chan []byte
}

var clientList = map[int32]*Client{}

func TCPStressTest(httpServer, tcpServer string, minInterval, maxInterval float32) {
	defer RecoverPanic()
	argNum := len(os.Args)
	if argNum != 3 {
		fmt.Printf("Error args num: %v \nRight is: go run stress_test_tcp.go {BEGIN_UID} {CLIENT_NUM}\n", argNum)
		return
	}
	idx, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: parse begin uid error, err=%v\n", err)
		return
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: parse num error, err=%v\n", err)
		return
	}

	//init app space
	initAppSpace(httpServer)

	//设置随机数种子
	rand.Seed(int64(time.Now().Nanosecond()))

	//启动客户端
	for i := 0; i < num; i++ {
		beginClient(int32(idx+i), tcpServer, minInterval, maxInterval)
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			num := 0
			for range clientList {
				num++
			}
			fmt.Println("current client num: ", num)
		}
	}()
}

func beginClient(uid int32, addr string, minInterval, maxInterval float32) {
	conn, err := net.Dial("tcp", addr)
	check(err, "dial server failed:")
	wCh := make(chan []byte, 100)
	go clientWriter(wCh, conn)
	rCh := make(chan []byte, 100)
	go clientReader(rCh, conn)

	//正常登陆
	wCh <- login(T_APP_ID, uid)
	rsp := getMsg(rCh)
	unpack(pb.CmdLoginAck, rsp)
	//join space
	x, y := randPos()
	wCh <- join(T_SPACE_ID, x, y)
	rsp = getMsg(rCh)
	unpack(pb.CmdJoinAck, rsp)
	//user list
	checkRspMsg(rCh, pb.CmdUserListNtf)

	go clientTimer(uid, &Client{W: wCh, R: rCh}, minInterval, maxInterval)
}

func endClient(uid int32) {
	c, ok := clientList[uid]
	if ok {
		closeClient(c.W)
		delete(clientList, uid)
	}
}

func clientReader(rCh chan []byte, conn net.Conn) {
	defer RecoverPanic()
	reader(rCh, conn)
}
func clientWriter(wCh chan []byte, conn net.Conn) {
	defer RecoverPanic()
	writer(wCh, conn)
}
func clientTimer(uid int32, c *Client, minInterval, maxInterval float32) {
	defer RecoverPanic()
	defer endClient(uid)
	clientList[uid] = c
	for {
		t := minInterval + rand.Float32()*(maxInterval-minInterval)
		time.Sleep(time.Millisecond * time.Duration(int(t*1000)))
		//select {
		//case <-time.After(time.Millisecond * time.Duration(int(t*1000))):
		//case <-util.GlobalDie:
		//	return
		//}
		ackCmd, data := randMsg()
		c.W <- data
		checkRspMsg(c.R, ackCmd)
	}
}

//======================================================
//======================================================
//======================================================

func randPos() (float32, float32) {
	x := rand.Float32() * 100 //[0.0,1.0)
	y := rand.Float32() * 100 //[0.0,1.0)
	return x, y
}

// 返回   ack消息id，req数据包
func randMsg() (int32, []byte) {
	n := rand.Intn(10) //[0,10)
	switch n {
	case 0: // 10% broadcast
		return pb.CmdBroadcastAck, broadcast()
	default: // 90% move
		x, y := randPos()
		return pb.CmdMoveAck, move(0, x, y)
	}
}
