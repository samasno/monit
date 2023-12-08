package mock

import (
	"fmt"
	"io"
	"net"
	"sync"
)

var connId = 0

func VanillaTcpServer(addr string, wg *sync.WaitGroup) func() error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				println("Server closed, shutting down")
				break
			}
			connId++
			go func(c net.Conn, id int) {
				defer wg.Done()
				defer conn.Close()
				b := make([]byte, 1024)
				println(fmt.Sprintf("Accepted new connection id: %d\n", id))
				for {
					n, err := conn.Read(b)
					if err == io.EOF {
						println("End of stream connection id", id)
						break
					}
					if err != nil {
						println("handle connection error:", err.Error())
						break
					}
					if n > 0 {
						fmt.Printf("Read %d bytes from connection %d\n", n, id)
						println(string(b))
					}
				}
				println("closing handle connection for id", id)
			}(conn, connId)
		}
	}(ln)
	return func() error {
		wg.Wait()
		return ln.Close()
	}
}
