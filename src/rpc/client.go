package rpc

import (
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

// TcpClient : tcp client
type TcpClient struct {
	address   string
	protocol  ITransProtocol
	invokeNum int32
	reqPool   chan *ReqMsg
}

// worker conf
var (
	MaxReqCap    = 10000
	MaxWorkerNum = 50
)

// NewTcpClient :
func NewTcpClient(address string) *TcpClient {
	t := &TcpClient{address, &MyProtocol{}, 0, make(chan *ReqMsg, MaxReqCap)}
	for i := 0; i < MaxWorkerNum; i++ {
		go t.Work()
	}
	return t
}

// Invoke : rpc invoke
func (c *TcpClient) Invoke(svrName string, fucName string, reqBuf []byte,
	timeout int) (ret int, rspBuf []byte) {
	reqPkg := &RequestPkg{
		RequestID: atomic.AddInt32(&c.invokeNum, 1),
		SvrName:   svrName,
		FuncName:  fucName,
		ReqBuf:    reqBuf,
		// ReqTime:   time.Now().Unix(),
	}

	reqMsg := &ReqMsg{reqPkg, make(chan *ResponsePkg, 1)}
	go func() {
		c.reqPool <- reqMsg
	}()
	timeoutCh := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		timeoutCh <- true
	}()
	rspPkg := new(ResponsePkg)
	select {
	case rspPkg = <-reqMsg.Rsp:
	case <-timeoutCh:
		rspPkg.RetCode = PkgTimeout
	}
	return int(rspPkg.RetCode), rspPkg.RspBuf
}

// NewConnection : new connection
func (c *TcpClient) NewConnection() (*Connection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.address)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn.SetKeepAlive(true)
	// log.Println("new conn")
	return &Connection{conn, c.protocol}, nil
}

// Work : worker handle req
func (c *TcpClient) Work() {
	var conn *Connection
	var err error
	count := 2
	for reqMsg := range c.reqPool {
		// use long connection, if conn is closed by remote, so reconnect
		for num := 0; num < count; num++ {
			if conn == nil {
				conn, err = c.NewConnection()
				if err != nil {
					rsp := new(ResponsePkg)
					rsp.RetCode = ConCreateErr
					reqMsg.Rsp <- rsp
					log.Println(err)
					break
				}
			}
			rsp, err := c.invoke(reqMsg, conn)
			if err == nil { // normal return
				reqMsg.Rsp <- rsp
				break
			}
			if err == io.EOF { // closed by remote, need reconect
				conn.Close()
				conn = nil
				continue
			} else { // other err
				reqMsg.Rsp <- rsp
				break
			}
		}
	}
}

func (c *TcpClient) invoke(reqMsg *ReqMsg, conn *Connection) (*ResponsePkg, error) {
	rsp := new(ResponsePkg)
	sendData, _ := Encode(reqMsg.Req)
	_, err := conn.Send(sendData)
	if err != nil {
		rsp.RetCode = ConSendErr
		return rsp, err
	}

	recvBuf, err := conn.Recv(0)
	if err != nil {
		rsp.RetCode = ConRecvErr
		log.Println(err)
		return rsp, err
	}
	rsp, err = DecodeRsp(recvBuf)
	return rsp, err
}
