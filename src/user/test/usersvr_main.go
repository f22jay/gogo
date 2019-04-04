package main

import (
	"log"
	"rpc"
	"user"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	userSvr := &user.UserSvr{}
	tcpSvr := rpc.NewTcpServer(":12345", userSvr)
	if err := tcpSvr.Run(); err != nil {
		log.Fatal(err)
	}
}
