package mock

import (
	"net"
	"sync"
	"testing"
)

func TestTcpServer(t *testing.T) {
	wg := &sync.WaitGroup{}
	addr := ":8080"
	closer := VanillaTcpServer(addr, wg)
	if closer == nil {
		t.Fatal("Failed to launch server")
	}
	wg.Add(1)
	client1, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err.Error())
	}
	client1.Write([]byte("message 1 from test."))
	client1.Write([]byte("message 2 from test."))
	client1.Close()

	wg.Add(1)
	client2, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err.Error())
	}
	client2.Write([]byte("message 3 from test."))
	client2.Write([]byte("message 4 from test."))
	client2.Close()

	err = closer()
	if err != nil {
		t.Fatal(err.Error())
	}
}
