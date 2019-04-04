package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

// RequestPkg request pack
type RequestPkg struct {
	RequestID int32
	SvrName   string
	FuncName  string
	ReqBuf    []byte
	ReqTime   int64
}

// ResponsePkg response pack
type ResponsePkg struct {
	RequestID int32
	SvrName   string
	FuncName  string
	RspBuf    []byte
	RetCode   int32
}

// ReqMsg req sturct for client
type ReqMsg struct {
	Req *RequestPkg
	Rsp chan *ResponsePkg
}

// PkgDecodeFail decode data err
const (
	PkgEncodeErr = 1000
	PkgDecodeErr = 1001
	PkgTimeout   = 1002
	ConCreateErr = 1003
	ConSendErr   = 1004
	ConRecvErr   = 1005
)

// TypeError TypeError
type TypeError struct{}

func (e *TypeError) Error() string {
	return "no matched type"
}

// Encode request & response
func Encode(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if req, ok := obj.(*RequestPkg); ok {
		err := enc.Encode(req)
		if err != nil {
			return nil, err
		}

	} else if rsp, ok := obj.(*ResponsePkg); ok {
		err := enc.Encode(rsp)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, &TypeError{}
	}
	return buf.Bytes(), nil
}

// DecodeReq  decode req from bytes
func DecodeReq(data []byte) (*RequestPkg, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var req RequestPkg
	err := dec.Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// DecodeRsp decode rsp from byte
func DecodeRsp(data []byte) (*ResponsePkg, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var rsp ResponsePkg
	err := dec.Decode(&rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

// judge recv pkg
const (
	PkgFull = iota
	PkgLess
	PkgErr
)

const maxPkgLen = 10 * 1024 * 1024

//ITransProtocol network data protocol
type ITransProtocol interface {
	MakeProtocol([]byte) []byte
	ParseProtocol([]byte) (int, []byte)
}

// MyProtocol : my protocol
type MyProtocol struct{}

// MakeProtocol : make protocol data
func (p *MyProtocol) MakeProtocol(reqData []byte) []byte {
	res := make([]byte, 4)
	res = append(res, reqData...)
	binary.BigEndian.PutUint32(res[:4], uint32(len(res)))
	return res
}

// ParseProtocol parse net data
func (p *MyProtocol) ParseProtocol(netData []byte) (int, []byte) {
	realLen := uint32(len(netData))
	if realLen < 4 {
		return PkgLess, nil
	}
	length := binary.BigEndian.Uint32(netData[:4])
	if length > maxPkgLen || length < realLen {
		return PkgErr, nil
	}
	if length > realLen {
		return PkgLess, nil
	}
	realPkg := make([]byte, 0)
	realPkg = append(realPkg, netData[4:]...)
	return PkgFull, realPkg
}
