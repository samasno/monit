package forwarder

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"sync"
	"testing"

	mock "github.com/samasno/monit/pkg/agent/mocks"
	"github.com/samasno/monit/pkg/agent/types"
)

func TestTcpClientConnect(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.VanillaTcpServer(":8080", wg)
	upstream := &types.Upstream{
		Connection: nil,
		Url:        "localhost",
		Port:       8080,
		TlsConfig:  nil,
	}
	emitter := &mock.MockEmitter{}
	tcpForwarder := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  emitter,
	}
	err := tcpForwarder.Connect()
	if err != nil {
		println("Failed to connect")
		t.Fatal(err.Error())
	}
	tcpForwarder.Disconnect()
	closer()
}

func TestTcpClientDisconnect(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.VanillaTcpServer(":8080", wg)
	upstream := &types.Upstream{
		Connection: nil,
		Url:        "localhost",
		Port:       8080,
		TlsConfig:  nil,
	}
	emitter := &mock.MockEmitter{}
	tcpForwarder := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  emitter,
	}
	err := tcpForwarder.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = tcpForwarder.Disconnect()
	if err != nil {
		println("Failed to disconnect from upstream")
		t.Fatal(err.Error())
	}
	if tcpForwarder.Upstream.Connection != nil {
		t.Fatal("Upstream connection was not nilled")
	}
	closer()
}

func TestTcpClientPush(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.VanillaTcpServer(":8080", wg)
	upstream := &types.Upstream{
		Connection: nil,
		Url:        "localhost",
		Port:       8080,
		TlsConfig:  nil,
	}
	emitter := &mock.MockEmitter{}
	tcpForwarder := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  emitter,
	}
	err := tcpForwarder.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}
	msgOne := "message one."
	msgTwo := "message two."
	wg.Add(1)
	err = tcpForwarder.Push([]byte(msgOne))
	if err != nil {
		println("Failed to push first message")
		t.Error(err.Error())
	}
	err = tcpForwarder.Push([]byte(msgTwo))
	if err != nil {
		println("Failed to push first message")
		t.Error(err.Error())
	}
	tcpForwarder.Disconnect()
	closer()
}

func TestTlsClientConnect(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.TlsTcpServer(":8080", wg)
	if closer == nil {
		t.Fatal("Failed to start tls server")
	}
	cert, err := ioutil.ReadFile("../../../certs/testing.crt")
	if err != nil {
		t.Fatal(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	config := &tls.Config{RootCAs: caCertPool, ServerName: "locahost", InsecureSkipVerify: true}
	upstream := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: config,
	}
	fwd := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  &mock.MockEmitter{},
	}
	wg.Add(1)
	err = fwd.Connect()
	if err != nil {
		t.Fatal(err)
	}
	fwd.Disconnect()
	closer()
}

func TestTlsClientDisconnect(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.TlsTcpServer(":8080", wg)
	if closer == nil {
		t.Fatal("Failed to start tls server")
	}
	cert, err := ioutil.ReadFile("../../../certs/testing.crt")
	if err != nil {
		t.Fatal(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	config := &tls.Config{RootCAs: caCertPool, ServerName: "locahost", InsecureSkipVerify: true}
	upstream := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: config,
	}
	fwd := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  &mock.MockEmitter{},
	}
	wg.Add(1)
	err = fwd.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = fwd.Disconnect()
	if err != nil {
		t.Fatal(err.Error())
	}
	closer()
	println("Succesfully disconnected")
}

func TestTlsClientPush(t *testing.T) {
	wg := &sync.WaitGroup{}
	closer := mock.TlsTcpServer(":8080", wg)
	if closer == nil {
		t.Fatal("Failed to start tls server")
	}
	cert, err := ioutil.ReadFile("../../../certs/testing.crt")
	if err != nil {
		t.Fatal(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	config := &tls.Config{RootCAs: caCertPool, ServerName: "locahost", InsecureSkipVerify: true}
	upstream := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: config,
	}
	fwd := ForwarderTcpClient{
		Upstream: upstream,
		Emitter:  &mock.MockEmitter{},
	}
	wg.Add(1)
	err = fwd.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = fwd.Push([]byte("Test tls message one."))
	if err != nil {
		t.Fatal(err)
	}
	err = fwd.Push([]byte("Test tls message two with extra bytes."))
	if err != nil {
		t.Fatal(err)
	}
	err = fwd.Disconnect()
	if err != nil {
		t.Fatal(err.Error())
	}
	closer()
	println("Succesfully disconnected")
}
