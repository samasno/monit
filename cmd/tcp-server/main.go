package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var connId = 0

func main() {
	println("Starting mock tcp server at", ":8080")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err.Error())
	}
	println("Mock tcp server running at", ":8080")
	for {
		conn, err := ln.Accept()
		if err != nil {
			println("mocktcp: Server closed, shutting down")
			break
		}
		connId++
		go func(c net.Conn, id int) {
			defer conn.Close()
			println(fmt.Sprintf("mocktcp: Accepted new connection id: %d\n", id))
			for {
				b := make([]byte, 1024)
				n, err := conn.Read(b)
				if err == io.EOF {
					println("mocktcp: End of stream connection id", id)
					break
				}
				if err != nil {
					println("mocktcp: handle connection error:", err.Error())
					break
				}
				if n > 0 {
					fmt.Printf("mocktcp:Read %d bytes from connection %d\n", n, id)
					println("mocktcp:" + string(b[:n]))
				}
			}
			println("mocktcp: closing handle connection for id", id)
		}(conn, connId)
	}
}
