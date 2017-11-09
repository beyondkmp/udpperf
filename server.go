package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
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

	quit := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go listen(conn, quit)
	}
	<-quit // hang until an error
}

func listen(connection *net.UDPConn, quit chan struct{}) {
	buffer := make([]byte, 1024)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for {
		n, remoteAddr, err = connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(n, remoteAddr, err)
			break
		}
	}
	fmt.Println("listener failed - ", err)
	quit <- struct{}{}
}
