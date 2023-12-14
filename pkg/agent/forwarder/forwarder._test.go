package forwarder

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"
	"testing"
	"time"

	mock "github.com/samasno/monit/pkg/agent/mocks"
	"github.com/samasno/monit/pkg/agent/types"
)

func TestForwarderConnectCloseTcp(t *testing.T) {
	testsock := "./test.sock"
	us := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: nil,
	}
	tcpFwd := &ForwarderTcpClient{
		Upstream: us,
		Emitter:  &mock.MockEmitter{},
	}
	ds := &types.Downstream{
		Url: testsock,
	}
	dgLn := &UnixDatagramSocketListener{
		Downstream: ds,
		Logger:     &mock.MockEmitter{},
	}
	fwd := Forwarder{
		UpstreamClient:     tcpFwd,
		DownstreamListener: dgLn,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	closer := mock.VanillaTcpServer(":8080", wg)
	err := fwd.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fwd.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer closer()
}

func TestForwarderConnectCloseTlsTcp(t *testing.T) {
	cert, err := ioutil.ReadFile("../../../certs/testing.crt")
	if err != nil {
		t.Fatal(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{RootCAs: caCertPool, ServerName: "localhost", InsecureSkipVerify: true}
	testsock := "./test.sock"
	us := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: tlsConfig,
	}
	tcpFwd := &ForwarderTcpClient{
		Upstream: us,
		Emitter:  &mock.MockEmitter{},
	}
	ds := &types.Downstream{
		Url: testsock,
	}
	dgLn := &UnixDatagramSocketListener{
		Downstream: ds,
		Logger:     &mock.MockEmitter{},
	}
	fwd := Forwarder{
		UpstreamClient:     tcpFwd,
		DownstreamListener: dgLn,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	closer := mock.TlsTcpServer(":8080", wg)
	err = fwd.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fwd.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer closer()
}

func TestForwarderRunPushTcp(t *testing.T) {
	testsock := "./test.sock"
	us := &types.Upstream{
		Url:       "localhost",
		Port:      8080,
		TlsConfig: nil,
	}
	tcpFwd := &ForwarderTcpClient{
		Upstream: us,
		Emitter:  &mock.MockEmitter{},
	}
	ds := &types.Downstream{
		Url: testsock,
	}
	dgLn := &UnixDatagramSocketListener{
		Downstream: ds,
		Logger:     &mock.MockEmitter{},
	}
	fwd := Forwarder{
		UpstreamClient:     tcpFwd,
		DownstreamListener: dgLn,
	}
	wg := &sync.WaitGroup{}
	closer := mock.VanillaTcpServer(":8080", wg)
	err := fwd.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	if err != nil {
		t.Fatal(err.Error())
	}
	wg.Add(1)
	client, err := net.Dial("unixgram", testsock)
	if err != nil {
		t.Fatal(err.Error())
	}
	client.Write([]byte("Test client one"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	clientTwo, err := net.Dial("unixgram", testsock)
	if err != nil {
		t.Fatal(err.Error())
	}
	clientTwo.Write([]byte("Test client one"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = clientTwo.Close()
	if err != nil {
		t.Fatal(err.Error())
	}

	time.Sleep(1 * time.Second)
	fwd.Close()
	closer()
}
