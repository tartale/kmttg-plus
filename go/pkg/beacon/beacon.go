package main

import (
	"fmt"
	"net"
	"time"
)

const (
	MulticastIp   = "239.255.255.250"
	MulticastPort = 2190
	TimeoutMs     = 20000
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", MulticastIp, MulticastPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(TimeoutMs * time.Millisecond))

	for {
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				break
			}
			panic(err)
		}
		fmt.Printf("Received packet from %v: %s\n", addr, string(buf))
	}
}
