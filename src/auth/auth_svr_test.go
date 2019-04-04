package auth

import (
	"log"
	"strings"
	"testing"
)

func TestLoginOk(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	req := new(LoginByPassReq)
	rsp := new(LoginByPassRsp)
	req.UserAuth.Username = "jack"
	req.UserAuth.Passwd = "jack"
	svr := new(AuthSvr)
	ret := svr.LoginByPass(req, rsp)
	if ret != 0 {
		t.Error(req, " login failed ret", ret)
	}
	t.Log(req, rsp)
	tReq := new(LoginByTokenReq)
	tReq.Token = rsp.Token
	tRsp := new(LoginByTokenRsp)

	ret = svr.LoginByToken(tReq, tRsp)
	if ret != 0 {
		t.Error(tReq, " login failed ", ret)
	}
	if strings.Compare(tRsp.Username, req.UserAuth.Username) != 0 {
		t.Error(tRsp, req)
	}
	t.Log(tReq, tRsp)

}

func TestLoginErr(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	req := new(LoginByPassReq)
	rsp := new(LoginByPassRsp)
	req.UserAuth.Username = "jack"
	req.UserAuth.Passwd = "i forget passwd"
	svr := new(AuthSvr)
	ret := svr.LoginByPass(req, rsp)
	if ret == 0 {
		t.Error(req, " login suc ret", ret)
	}
	t.Log(req, rsp)

}

// insert test data
// func TestInsert(t *testing.T) {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	for begin := 0; begin < 10000000; begin += 10000 {
// 		if err := InsertData(begin); err != nil {
// 			t.Error(err)
// 		}
// 	}

// }
