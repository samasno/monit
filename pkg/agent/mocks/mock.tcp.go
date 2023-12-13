package mock

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

var connId = 0
var tlsConnId = 0

func VanillaTcpServer(addr string, wg *sync.WaitGroup) func() error {
	println("Starting mock tcp server at", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil
	}
	println("Mock tcp server running at", addr)
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				println("mocktcp: Server closed, shutting down")
				break
			}
			connId++
			go func(c net.Conn, id int) {
				defer wg.Done()
				defer conn.Close()
				b := make([]byte, 1024)
				println(fmt.Sprintf("mocktcp: Accepted new connection id: %d\n", id))
				for {
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
	}(ln)
	return func() error {
		wg.Wait()
		return ln.Close()
	}
}

func TlsTcpServer(addr string, wg *sync.WaitGroup) func() error {
	cert, err := tls.LoadX509KeyPair("../../../certs/testing.crt", "../../../certs/testing.pem")
	if err != nil {
		println("mocktlstcp:failed to load key pair")
		println(err.Error())
		return nil
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", addr, config)
	if err != nil {
		println("mocktlstcp: Failed to start tls server")
		println(err.Error())
		return nil
	}
	go func(ln net.Listener, wg *sync.WaitGroup) {
		println("mocktlstcp: server listening at " + addr)
		for {
			conn, err := ln.Accept()
			if err != nil {
				println("mocktlstcp: Server closed, shutting down")
				break
			}
			tlsConnId++
			go func(conn net.Conn, id int) {
				defer conn.Close()
				defer wg.Done()
				b := make([]byte, 1024)
				for {
					n, err := conn.Read(b)
					if err == io.EOF {
						println("mocktlstcp: Reached end of stream tls id: " + strconv.Itoa(id))
						break
					}
					if err != nil {
						println(err.Error())
						break
					}
					if n > 0 {
						fmt.Printf("mocktlstcp: Read %d bytes from connection %d\n", n, id)
						println(string(b[:n]))
					}
				}
				println("mocktlstcp: Closing connection handler for " + strconv.Itoa(id))
			}(conn, tlsConnId)
		}
	}(ln, wg)
	return func() error {
		println("mocktlstcp: Draining connections")
		wg.Wait()
		println("mocktlstcp: No connections remaining")
		println("mocktlstcp: Closing listener")
		return ln.Close()
	}
}
