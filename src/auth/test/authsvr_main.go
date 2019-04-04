package main

import (
	"auth"
	"log"
	"rpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	authSvr := &auth.AuthSvr{}
	tcpSvr := rpc.NewTcpServer(":23456", authSvr)
	if err := tcpSvr.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("over")
}
