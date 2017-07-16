package rpc

import (
	"golang.org/x/net/context"
	"net"
	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/SongCF/scene/model"
	"github.com/SongCF/scene/pb"
	"github.com/SongCF/scene/util"
)

type server struct {
}


func (s *server) Notice(ctx context.Context, req *NoticeRequest) (*NoticeReply, error) {
	util.RecoverPanic()
	log.Println("handle rpc call ...")
	ret := &NoticeReply{}
	sess := model.GetSession(string(req.AppId), req.UserId)
	if sess == nil {
		ret.Code = pb.ErrUserOffline.Id
		return ret, nil
	} else {
		log.Println("rpc -> handle call: cmd=", req.Cmd)
		sess.Rsp2(req.Cmd, req.Payload)
		return ret, nil
	}
}


func InitServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterUserServer(s, &server{})
	s.Serve(lis)
}