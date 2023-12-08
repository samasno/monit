package mock

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
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

func TestTlsTcpServer(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := TlsTcpServer(":8080", wg)
	if closer == nil {
		t.Fatal("Failed to start server")
	}
	cert, err := ioutil.ReadFile("../../../certs/testing.crt")
	if err != nil {
		t.Fatal(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{RootCAs: caCertPool, ServerName: "localhost", InsecureSkipVerify: true}
	client, err := tls.Dial("tcp", ":8080", tlsConfig)
	if err != nil {
		t.Fatal(err.Error())
	}
	wg.Add(1)
	client.Write([]byte("tls test message"))
	client.Close()
	err = closer()
	if err != nil {
		println("Failed to stop server")
		t.Fatal(err.Error())
	}
}
