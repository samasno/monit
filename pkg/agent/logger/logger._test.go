package logger

import (
	"encoding/json"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/samasno/monit/pkg/agent/listeners"
	mock "github.com/samasno/monit/pkg/agent/mocks"
	"github.com/samasno/monit/pkg/agent/types"
	"github.com/samasno/monit/pkg/agent/vars"
)

var testSocket = "./test.sock"
var testClientSocket = "./test-client.sock"
var testLogFile = "./test.log"

func TestOpenCloseLogger(t *testing.T) {
	downstream := &types.Downstream{
		Url: testSocket,
	}
	testListener := &listeners.UnixDatagramSocketListener{
		Name:       "unix-datagram-listener",
		Downstream: downstream,
		Logger:     &mock.MockEmitter{},
	}
	logger := Logger{
		Listener: testListener,
		LogFile:  testLogFile,
	}
	err := logger.connect()
	if err != nil {
		println("Failed to open")
		t.Fatal(err.Error())
	}
	err = logger.close()
	if err != nil {
		println("Failed to close")
		t.Fatal(err.Error())
	}
}

func TestListenAndLog(t *testing.T) {
	downstream := &types.Downstream{
		Url: testSocket,
	}
	testListener := &listeners.UnixDatagramSocketListener{
		Name:       "unix-datagram-listener",
		Downstream: downstream,
		Logger:     &mock.MockEmitter{},
	}
	logger := Logger{
		Listener: testListener,
		LogFile:  testLogFile,
	}
	// listen and log
	err := logger.ListenAndLog()
	if err != nil {
		t.Fatal(err.Error())
	}
	// create event
	eventOne := types.Event{
		Type: vars.FORWARDER_CLIENT_LOG,
		Payload: types.Payload{
			Source:  "Logger-testing",
			Message: "Sent test message one",
			Level:   vars.INFO,
		},
	}
	eventTwo := types.Event{
		Type: vars.FORWARDER_CLIENT_LOG,
		Payload: types.Payload{
			Source:  "Logger-testing",
			Message: "Sent test message two",
			Level:   vars.INFO,
		},
	}
	eventOneJson, err := json.Marshal(eventOne)
	if err != nil {
		t.Fatal(err.Error())
	}
	eventTwoJson, err := json.Marshal(eventTwo)
	if err != nil {
		t.Fatal(err.Error())
	}
	raddr, err := net.ResolveUnixAddr("unixgram", testSocket)
	if err != nil {
		t.Fatal(err.Error())
	}
	laddr, err := net.ResolveUnixAddr("unixgram", testClientSocket)
	if err != nil {
		t.Fatal(err.Error())
	}
	client, err := net.DialUnix("unixgram", laddr, raddr)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = client.Write(eventOneJson)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = client.Write(eventTwoJson)
	if err != nil {
		t.Fatal(err.Error())
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(3 * time.Second)
		logger.close()
		wg.Done()
	}(wg)
	wg.Wait()
	println()
	println()
	logOutput, err := os.ReadFile(testLogFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	println("Got output in test log file")
	println(string(logOutput))
	os.RemoveAll(testClientSocket)
	os.RemoveAll(testLogFile)
}
