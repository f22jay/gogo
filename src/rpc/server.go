package rpc

import (
	"log"
	"net"
)

// TcpServer send & recv msg, and call dispatcher
type TcpServer struct {
	listener     *net.TCPListener
	connNum      int32
	address      string
	invokeNum    int32
	isRunning    bool
	protocol     ITransProtocol
	keeplivetime int // seconds
	servant      Servant
}

// Servant : logic svr
type Servant interface {
	Dispatcher(req *RequestPkg) (rsp *ResponsePkg)
}

// NewTcpServer rpc svr
func NewTcpServer(address string, servant Servant) *TcpServer {
	t := &TcpServer{}
	t.address = address
	t.protocol = &MyProtocol{}
	t.keeplivetime = 60
	t.servant = servant
	t.isRunning = true
	return t
}

// Stop : stop server
func (t *TcpServer) Stop() {
	t.isRunning = false
}

func (t *TcpServer) listen() (err error) {
	addr, err := net.ResolveTCPAddr("tcp4", t.address)
	if err != nil {
		return err
	}
	t.listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		return err
	}
	log.Println("listen: ", t.address)
	return nil
}

func (t *TcpServer) handleConn(c *Connection) {
	defer c.Close()
	for {
		// every long connecton period, if not read, for over & close connection
		buf, err := c.Recv(t.keeplivetime)
		if err != nil {
			// log.Println(err)
			if _, ok := err.(PackageErr); !ok { // not logic err, net err
				return
			}
			continue
		}
		req, err := DecodeReq(buf)
		if err != nil {
			log.Println("decode err ", err)
			continue
		}
		rsp := t.servant.Dispatcher(req)
		buf, err = Encode(rsp)
		c.Send(buf)
	}
}

// Run server
func (t *TcpServer) Run() (err error) {
	if err := t.listen(); err != nil {
		return nil
	}

	for t.isRunning {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			log.Println("accept error")
			return err
		}
		// log.Println("get connection remote addr ", conn.RemoteAddr())
		conn.SetKeepAlive(true)
		c := &Connection{conn, t.protocol}
		go t.handleConn(c)
	}
	return nil
}
