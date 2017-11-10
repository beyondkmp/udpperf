package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "37", "port")
var numClients = flag.Int("num", 5, "numClients")
var sizePacket = flag.Int("size", 5, "sizePacket")
var secondSleep = flag.Int("second", 1, "secondSleep")
var wg sync.WaitGroup

var size int64

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	content := make([]byte, *sizePacket)
	rand.Read(content)
	for i := 0; i < *numClients; i++ {
		wg.Add(1)
		go func() {
			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				fmt.Println("Can't dial: ", err)
				os.Exit(1)
			}
			defer func() {
				wg.Done()
				conn.Close()
			}()

			timeoutChan := make(chan bool, 1)
			go func() {
				time.Sleep(time.Second * time.Duration(*secondSleep))
				timeoutChan <- true
			}()

			for {
				select {
				case <-timeoutChan:
					fmt.Println("finish")
					break
				default:
					_, err = conn.Write(content)
					if err != nil {
						fmt.Println("failed:", err)
						break
					}
					atomic.AddInt64(&size, int64(*sizePacket))
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("total send %v B\n", atomic.LoadInt64(&size))
}
