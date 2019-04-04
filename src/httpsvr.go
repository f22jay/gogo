package main

import (
	"auth"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"user"
)

// HttpRsp : general http response
type HttpRsp struct {
	Ret   int
	Msg   string
	Value string
}

// err code
const (
	ErrArg   = 3001
	ErrLogin = 3002
)

var random *rand.Rand

func newHttpRsp(ret int, msg string, value string) string {
	rsp := &HttpRsp{ret, msg, value}
	data, _ := json.Marshal(&rsp)
	return string(data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("parseform err:", err)
		io.WriteString(w, newHttpRsp(ErrArg, "get data err", ""))
		return
	}
	client := auth.NewAuthClient("127.0.0.1:23456", 3)
	var req auth.LoginByPassReq
	var rsp auth.LoginByPassRsp
	req.UserAuth.Username = r.PostFormValue("username")
	req.UserAuth.Passwd = r.PostFormValue("passwd")
	ret := client.LoginByPass(&req, &rsp)
	if ret != 0 {
		log.Println("login fail ret:", ret)
		io.WriteString(w, newHttpRsp(ErrLogin, "login failed", ""))
		return
	}
	data, _ := json.Marshal(&rsp)
	io.WriteString(w, newHttpRsp(0, "", string(data)))
}

func loginByToken(token string) (int, string) {
	client := auth.NewAuthClient("127.0.0.1:23456", 3)
	var req auth.LoginByTokenReq
	var rsp auth.LoginByTokenRsp
	req.Token = token
	ret := client.LoginByToken(&req, &rsp)
	if ret != 0 {
		log.Println(req, " login failed")
		return ret, ""
	}
	return 0, rsp.Username
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// 1. check login & get username
	cookie, err := r.Cookie("token")
	if err != nil {
		io.WriteString(w, newHttpRsp(ErrArg, "no token", ""))
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("get data err ", err)
		io.WriteString(w, newHttpRsp(ErrArg, "get data err", ""))
		return
	}

	ret, username := loginByToken(cookie.Value)
	if ret != 0 {
		io.WriteString(w, newHttpRsp(ErrLogin, "login failed", ""))
		log.Println("login failed ret ", ret)
		return
	}

	client := user.NewUserClient("127.0.0.1:12345", 3)
	// 2. get or upload
	if r.Method == "GET" {
		var req user.GetByNameReq
		var rsp user.GetByNameRsp
		if username == "test" {
			req.Username = "test" + strconv.Itoa(random.Intn(9999999))
		} else {
			req.Username = username
		}

		ret := client.GetByName(&req, &rsp)
		if ret != 0 {
			log.Println(ret)
		}
		if req.Username[0:4] == "test" && req.Username != rsp.UserItem.Username {
			log.Println(req, rsp)
		}
		data, _ := json.Marshal(&rsp)
		io.WriteString(w, newHttpRsp(ret, "", string(data)))
		return
	}

	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm != nil {
			if len(r.MultipartForm.Value["username"]) == 0 ||
				len(r.MultipartForm.Value["nickname"]) == 0 {
				io.WriteString(w, newHttpRsp(ErrArg, "no username or nickname", ""))
				return
			}
			if r.MultipartForm.Value["username"][0] != username {
				io.WriteString(w, newHttpRsp(ErrArg, " username not matched", ""))
				return
			}

		}
		_, h, err := r.FormFile("image")
		if err != nil {
			io.WriteString(w, newHttpRsp(ErrArg, "no upload file", ""))
			return
		}
		var ureq user.UpdateByNameReq
		var ursp user.UpdateByNameRsp
		ureq.UserItem.Username = username
		ureq.UserItem.Nickname = r.MultipartForm.Value["nickname"][0]
		ureq.UserItem.Profile = h.Filename
		ret := client.UpdateByName(&ureq, &ursp)
		if ret != 0 {
			log.Println(ret)
			io.WriteString(w, newHttpRsp(ret, "update failed", ""))
			return
		}
		io.WriteString(w, newHttpRsp(0, "", ""))
	}
}

func main() {
	random = rand.New(rand.NewSource(time.Now().Unix()))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("begin http")
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/login", loginHandler)
	err := http.ListenAndServe(":10086", nil)
	if err != nil {
		log.Fatal("listenAndServer: ", err.Error())
	}
}
