package main

import (
	"log"
	// "strings"
	"sync/atomic"
	"user"
)

var suc int32
var fail int32

func testClient(ch chan int) {
	client := user.NewUserClient("127.0.0.1:12345", 3)
	var req user.GetByNameReq
	var rsp user.GetByNameRsp
	req.Username = "hello"
	for i := 0; i < 500; i++ {

		ret := client.GetByName(&req, &rsp)
		if ret != 0 {
			atomic.AddInt32(&fail, 1)
		} else {
			atomic.AddInt32(&suc, 1)
		}
		// log.Println(req, rsp)
	}

	ch <- 1
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	num := 100
	count := 0
	ch := make(chan int, num)
	for count < num {
		count++
		go testClient(ch)
	}

	for v := range ch {
		num--
		if num < 2 {
			break
		}
		if v != 1 {
			log.Fatal("err")
		}
	}
	log.Println("suc ", suc, " fail ", fail)

}
