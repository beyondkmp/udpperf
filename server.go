package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "37", "port")
var size int64

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Printf("total receive %v B\n", atomic.LoadInt64(&size))
		}
	}()

	quit := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go listen(conn, quit)
	}
	<-quit // hang until an error
}

func listen(connection *net.UDPConn, quit chan struct{}) {
	buffer := make([]byte, 1500)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for {
		n, remoteAddr, err = connection.ReadFromUDP(buffer)
		atomic.AddInt64(&size, int64(n))
		if err != nil {
			fmt.Println(n, remoteAddr, err)
			break
		}
	}
	fmt.Println("listener failed - ", err)
	quit <- struct{}{}
}
