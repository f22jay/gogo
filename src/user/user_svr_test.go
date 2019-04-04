package user

import (
	"log"
	"testing"
)

func TestUpdate(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	req := new(GetByNameReq)
	rsp := new(GetByNameRsp)
	req.Username = "jack"
	svr := new(UserSvr)
	if ret := svr.GetByName(req, rsp); ret != 0 {
		t.Error(req, " get err ret:", ret)
		return
	}

	uReq := new(UpdateByNameReq)
	uReq.UserItem.Username = req.Username
	uReq.UserItem.Nickname = "nnnn"
	uReq.UserItem.Profile = "pppp"
	uRsp := new(UpdateByNameRsp)

	if ret := svr.UpdateByName(uReq, uRsp); ret != 0 {
		t.Error(uReq, " update err ret:", ret)
		return
	}

	if ret := svr.GetByName(req, rsp); ret != 0 {
		t.Error(req, " get err ret:", ret)
		return
	}

	if rsp.UserItem.Nickname != uReq.UserItem.Nickname || rsp.UserItem.Profile != uReq.UserItem.Profile {
		t.Error("update before:", uReq, " after: ", rsp)
	}

}
