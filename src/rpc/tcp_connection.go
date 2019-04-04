package rpc

import (
	"log"
	"net"
	"time"
)

// Connection : a rpc connection, just recv & send, use protocl
type Connection struct {
	conn     *net.TCPConn
	protocol ITransProtocol
}

// PackageErr : for parse package err
type PackageErr struct{}

func (PackageErr) Error() string {
	return "package parse err"
}

// Recv data
func (c *Connection) Recv(ts int) ([]byte, error) {
	buf := make([]byte, 4*1024)
	var recvBuf []byte
	for {
		if ts != 0 {
			c.conn.SetReadDeadline(time.Now().Add(time.Duration(ts) * time.Second))
		}
		n, err := c.conn.Read(buf)
		if err != nil {
			return nil, err
		}
		recvBuf = append(recvBuf, buf[:n]...)
		ret, dataBuf := c.protocol.ParseProtocol(recvBuf)
		if ret == PkgFull {
			return dataBuf, nil
		}
		if ret == PkgLess {
			continue
		}
		log.Println("err buf len: ", len(recvBuf))
		return nil, PackageErr{}
	}
}

// Send data
func (c *Connection) Send(data []byte) (int, error) {
	sendbuf := c.protocol.MakeProtocol(data)
	_, err := c.conn.Write(sendbuf)
	if err != nil {
		// log.Println("send err ", err)
		return 0, err
	}
	// log.Println("send buf len ", len(sendbuf))
	return len(data), nil
}

// Close Connection
func (c *Connection) Close() {
	// log.Println("close conn")
	c.conn.Close()
}
