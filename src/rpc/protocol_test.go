package rpc

import (
	"testing"
)

func Test_Encode(t *testing.T) {
	req := &RequestPkg{1, "test", "hello", []byte{1, 2, 3, 4}, 43434}
	buf, err := Encode(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(buf)
	decodeReq, err := DecodeReq(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(*decodeReq)

	rsp := &ResponsePkg{1, "test", "sdfdsf", []byte{1, 23, 3, 4}, 0}
	buf, err = Encode(rsp)
	if err != nil {
		t.Error(err)
	}
	decodeRsp, err := DecodeRsp(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(*decodeRsp)
}

func TestProtocol(t *testing.T) {

	req := &RequestPkg{1, "test", "hello", []byte{1, 2, 3, 4}, 43434}
	buf, err := Encode(req)
	if err != nil {
		t.Error(err)
	}

	// Full case
	protocol := MyProtocol{}
	protocolBuf := protocol.MakeProtocol(buf)
	flag, parseBuf := protocol.ParseProtocol(protocolBuf)
	if flag != PkgFull {
		t.Error(flag)
	}
	if len(buf) != len(parseBuf) {
		t.Error("wrong len buf len: ", len(buf), " protocolBuf len: ", len(protocolBuf), "", len(parseBuf))
	}

	for idx, value := range buf {
		if value != parseBuf[idx] {
			t.Error("protocol error  ")
		}
	}

	// Less case
	wrongBuf := make([]byte, len(protocolBuf)/2)
	copy(wrongBuf, protocolBuf[0:len(protocolBuf)/2])

	flag, parseBuf = protocol.ParseProtocol(wrongBuf)
	if flag != PkgLess {
		t.Error("parse less case err flag ", flag)
	}
}
