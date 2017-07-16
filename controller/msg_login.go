package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	. "github.com/SongCF/scene/model"
	"github.com/SongCF/scene/pb"
)

//user login
func LoginReq(s *Session, m []byte) (int32, proto.Message) {
	payload := &pb.LoginReq{}
	err := proto.Unmarshal(m, payload)
	if err != nil {
		return pb.Error(pb.CmdLoginReq, pb.ErrMsgFormat)
	}
	// check the session already login
	if s.HasLogin() {
		return pb.Error(pb.CmdLoginReq, pb.ErrDuplicateLogin)
	}

	appId := string(payload.AppId)
	uid := payload.GetUserId()

	// 1.check has app
	if !HasApp(appId) {
		return pb.Error(pb.CmdLoginReq, pb.ErrAppNotExist)
	}
	// kick out other session login by the uid if had
	oldSess := GetSession(appId, uid)
	if oldSess != nil {
		//ntf
		t := OFFLINE_TYPE_OTHER_LOGIN
		offlineMsg := &pb.OfflineNtf{
			Type: &t,
		}
		log.Infof("kick user(%v:%v) old session", oldSess.AppId, oldSess.Uid)
		oldSess.Rsp(pb.CmdOfflineNtf, offlineMsg)
		//clean cache
		Logout(oldSess)
	}
	// set cache (init UserInfo in cache)
	conn, err := CCPool.Get()
	if err != nil {
		log.Errorf("LoginReq(user[%v:%v]) CCPool:Get error(%v)", appId, uid, err)
		return pb.Error(pb.CmdLoginReq, pb.ErrServerBusy)
	}
	defer CCPool.Put(conn)
	e := ResetUserInfo(appId, uid, conn)
	if e != nil {
		log.Errorf("LoginReq(user[%v:%v]) ResetUserInfo error(%v)", appId, uid, err)
		return pb.Error(pb.CmdLoginReq, pb.ErrServerBusy)
	}
	// set session user info
	s.AppId = appId
	s.Uid = uid
	SetSession(appId, uid, s)

	rsp := &pb.LoginAck{}
	return pb.CmdLoginAck, rsp
}

func Logout(s *Session) {
	if s != nil {
		// post leave
		Leave(s, &pb.LeaveReq{})
		// delete session
		DelSession(s.AppId, s.Uid)
		// clean session
		s.Clean()
	}
}
