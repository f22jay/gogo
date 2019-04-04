package main

import (
	"auth"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"sync/atomic"
	"time"
)

var suc int32
var fail int32

// var cost int64
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func testClient(ch chan int) {
	client := auth.NewAuthClient("127.0.0.1:23456", 3)
	var req auth.LoginByPassReq
	req.UserAuth.Username = "qps"
	req.UserAuth.Passwd = "1234234"
	expectToken := req.UserAuth.Username + req.UserAuth.Passwd
	// now := time.Now().UnixNano() / 1e6
	for i := 0; i < 10000; i++ {
		var rsp auth.LoginByPassRsp
		ret := client.LoginByPass(&req, &rsp)
		// end := time.Now().UnixNano() / 1e6
		// cs := end - now
		// atomic.AddInt64(&cost, cs)
		if ret != 0 || strings.Compare(rsp.Token, expectToken) != 0 {
			atomic.AddInt32(&fail, 1)
		} else {
			atomic.AddInt32(&suc, 1)
		}
		// log.Println(req, rsp)
		// now = end
	}

	ch <- 1
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	num := 30
	count := 0
	ch := make(chan int, num)
	now := time.Now().UnixNano() / 1e6
	for count < num {
		count++
		go testClient(ch)
	}

	for v := range ch {
		num--
		if num == 0 {
			break
		}
		if v != 1 {
			log.Fatal("err")
		}
	}
	end := time.Now().UnixNano() / 1e6
	log.Println("suc ", suc, " fail ", fail)
	log.Println("qps:", float64(suc+fail)*1000/float64(end-now))
}
