package _test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"jhqc.com/songcf/scene/pb"
	"net"
	"strings"
	"time"
)

func startClient(tcpAddr string) (chan []byte, chan []byte) {
	conn, err := net.Dial("tcp", tcpAddr)
	check(err, "dial server failed:")
	wCh := make(chan []byte)
	go writer(wCh, conn)
	rCh := make(chan []byte)
	go reader(rCh, conn)
	return wCh, rCh
}

func closeClient(wCh chan []byte) {
	close(wCh)
}

func getMsg(ch chan []byte) []byte {
	select {
	case m := <-ch:
		return m
	case <-time.After(time.Second * 5):
		err := errors.New("get msg timeout")
		check(err, "")
		return nil
	}
}
func getNilMsg(ch chan []byte) {
	select {
	case m := <-ch:
		err := errors.New(fmt.Sprintf("get unexpect msg:%v", m))
		check(err, "")
	case <-time.After(time.Second * 5):
	}
}
func checkRspMsg(ch chan []byte, cmd int32) {
	var m []byte
	for {
		select {
		case m = <-ch:
		case <-time.After(time.Second * 5):
			err := errors.New("get msg timeout")
			check(err, fmt.Sprintf("check rsp :%v,", pb.RCode[int(cmd)]))
			return
		}
		// parse
		packet := &pb.Packet{}
		err := proto.Unmarshal(m, packet)
		check(err, "getRspMsg unmarshal packet failed:")
		cmd2 := packet.GetCmd()
		if cmd2 == cmd {
			return
		}
	}
}

func writer(ch chan []byte, conn net.Conn) {
	defer conn.Close()
	const lenHead = 4
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		size := len(data)
		buf := make([]byte, size+lenHead)
		binary.BigEndian.PutUint32(buf[:lenHead], uint32(size))
		copy(buf[lenHead:], data)
		_, err := conn.Write(buf[:size+lenHead])
		check(err, "write msg failed:")
	}
}

func reader(ch chan []byte, conn net.Conn) {
	const lenHead = 4
	head := make([]byte, lenHead)
	for {
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		_, err := io.ReadFull(conn, head)
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, "read head failed:")
		size := int(binary.BigEndian.Uint32(head))
		// read data
		buf := make([]byte, size)
		_, err = io.ReadFull(conn, buf)
		if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
			break
		}
		check(err, "read pb failed:")
		ch <- buf[:size]
	}
}
