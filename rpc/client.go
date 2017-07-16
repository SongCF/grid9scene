package rpc

import (
	"google.golang.org/grpc"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"github.com/SongCF/scene/util"
)

func TestRPC(addr string, appId string, userId int32, cmd int32, payload []byte) *NoticeReply {
	util.RecoverPanic()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewUserClient(conn)

	req := &NoticeRequest{
		AppId:[]byte(appId),
		UserId:userId,
		Cmd:cmd,
		Payload:payload,
	}
	log.Println("begin rpc call ...")
	r, err := c.Notice(context.Background(), req)
	log.Println("end rpc call ...")
	if err != nil {
		log.Fatalf("rpc failed: %v", err)
	}
	log.Println("rpc -> r: ", r)
	return r
}