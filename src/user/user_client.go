package user

import (
	"bytes"
	"encoding/gob"
	"log"
	"rpc"
	"sync"
)

// UserClient rpc client
type UserClient struct {
	t       *rpc.TcpClient
	timeout int
}

var t *rpc.TcpClient
var once sync.Once

// NewUserClient : return userclient
func NewUserClient(address string, timeout int) *UserClient {
	once.Do(func() {
		t = rpc.NewTcpClient(address)
	})

	return &UserClient{t, timeout}
}

// GetByName get uesr by name
func (c *UserClient) GetByName(req *GetByNameReq, rsp *GetByNameRsp) int {
	// encode req
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(req)

	// call
	ret, recvBuf := t.Invoke("UserSvr", "GetByName", buf.Bytes(), c.timeout)
	if ret != 0 {
		return int(ret)
	}

	// decode rsp
	buf = new(bytes.Buffer)
	buf.Write(recvBuf)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(rsp)
	if err != nil {
		log.Println(err)
		return rpc.PkgDecodeErr
	}

	return 0
}

// UpdateByName  uesr by name
func (c *UserClient) UpdateByName(req *UpdateByNameReq, rsp *UpdateByNameRsp) int {
	// encode req
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(req)

	// call
	ret, recvBuf := t.Invoke("UserSvr", "UpdateByName", buf.Bytes(), c.timeout)
	if ret != 0 {
		return int(ret)
	}

	// decode rsp
	buf = new(bytes.Buffer)
	buf.Write(recvBuf)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(rsp)
	if err != nil {
		log.Println(err)
		return rpc.PkgDecodeErr
	}

	return 0
}
