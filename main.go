package main

import (
	"fmt"
	"httpfromtcp/request"
	"log"
	"net"
)

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	log.Printf("Listening on port %s\n", port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		req, err := request.RequestFromHeader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method:%s\n", req.RequestLine.Method)
		fmt.Printf("- Target:%s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version:%s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		req.Headers.ForEach(func(n, v string) {
			fmt.Printf("- %s: %s\n", n, v)
		})
	}
}
