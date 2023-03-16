//package main
//
//import (
//	"github.com/tartale/kmttg-plus/go/server/resolvers"
//	"log"
//	"net/http"
//	"os"
//
//	"github.com/99designs/gqlgen/graphql/handler"
//	"github.com/99designs/gqlgen/graphql/playground"
//)
//
//const defaultPort = "8080"
//
//func main() {
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = defaultPort
//	}
//
//	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &resolvers.Resolver{}}))
//
//	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
//	http.Handle("/query", srv)
//
//	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
//	log.Fatal(http.ListenAndServe(":"+port, nil))
//}

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
