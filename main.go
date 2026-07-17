package main

import (
	"fmt"
	"httpfromtcp/request"
	"log"
	"net"
)

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	out := make(chan string, 1)

// 	go func() {
// 		defer f.Close()
// 		defer close(out)

// 		str := ""
// 		for {
// 			data := make([]byte, 8)
// 			n, err := f.Read(data)
// 			if err != nil {
// 				break
// 			}

// 			data = data[:n]
// 			if i := bytes.IndexByte(data, '\n'); i != -1 {
// 				str += string(data[:i])
// 				data = data[i+1:]
// 				out <- str
// 				str = ""
// 			}
// 			str += string(data)
// 		}

// 		if len(str) != 0 {
// 			out <- str
// 		}
// 	}()

// 	return out
// }

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
	}
}
