package forwarder

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/samasno/monit/pkg/agent/types"
	"github.com/samasno/monit/pkg/agent/vars"
)

type ForwarderTcpClient struct {
	Upstream *types.Upstream
	Emitter  types.Emitter
	shutdown *sync.WaitGroup
}

func (t *ForwarderTcpClient) Connect(shutdown *sync.WaitGroup) error {
	if t.shutdown == nil {
		t.shutdown = shutdown
		t.shutdown.Add(1)
	}
	if t.Upstream.Connection != nil {
		t.log(vars.NOTICE, "Connection already exists")
		return nil
	}
	dest := fmt.Sprintf("%s:%d", t.Upstream.Url, t.Upstream.Port)
	var conn net.Conn
	var err error
	if t.Upstream.TlsConfig != nil {
		conn, err = tls.Dial("tcp", dest, t.Upstream.TlsConfig)
		if err != nil {
			msg := fmt.Sprintf("Failed to open tcp/tls connection to upstream %s: %s", dest, err.Error())
			t.log(vars.ERROR, msg)
			return fmt.Errorf(msg + "\n")
		}
	} else {
		conn, err = net.Dial("tcp", dest)
		if err != nil {
			msg := fmt.Sprintf("Failed to open tcp connection to upstream %s: %s", dest, err.Error())
			t.log(vars.ERROR, msg)
			return fmt.Errorf(msg + "\n")
		}
	}
	t.Upstream.Connection = conn
	t.log(vars.INFO, "Opened connection to upstream tcp server at "+dest)
	return nil
}

func (t *ForwarderTcpClient) Disconnect() error {
	if t.Upstream.Connection != nil {
		err := t.Upstream.Connection.Close()
		if err != nil {
			msg := fmt.Sprintf("Failed to disconnect from %s: %s", t.Upstream.Url, err.Error())
			t.log(vars.ERROR, msg)
			return fmt.Errorf(msg)
		}
		t.Upstream.Connection = nil
	}
	if t.shutdown != nil {
		t.shutdown.Done()
		t.shutdown = nil
	}
	t.log(vars.INFO, fmt.Sprintf("Disconnected from %s", t.Upstream.Url))
	return nil
}

func (t *ForwarderTcpClient) Push(payload []byte) error {
	if t.Upstream.Connection == nil {
		errmsg := "Connection to upstream is closed"
		t.log(vars.ERROR, errmsg)
		return fmt.Errorf(errmsg + "\n")
	}
	_, err := t.Upstream.Connection.Write(payload)
	if err != nil {
		errmsg := fmt.Sprintf("Failed to forward message to tcp upstream: %s", err.Error())
		t.log(vars.ERROR, errmsg)
		return fmt.Errorf(errmsg)
	}
	msg := fmt.Sprintf("Pushed %d bytes to %s", len(payload), t.Upstream.Url)
	t.log(vars.INFO, msg)
	return nil
}

func (t *ForwarderTcpClient) log(level int, message string) error {
	if t.Emitter == nil {
		return fmt.Errorf("No emitter to send logs")
	}
	payload := types.Payload{
		Source:  NAME,
		Message: message,
		Level:   level,
	}
	event := types.Event{
		Type:    vars.FORWARDER_CLIENT_LOG,
		Payload: payload,
	}
	err := t.Emitter.Emit(event)
	if err != nil {
		msg := fmt.Sprintf("Failed to emit message")
		return fmt.Errorf(msg + "\n")
	}
	return nil
}

func (t *ForwarderTcpClient) Ok() (bool, string) {
	ok := true
	msg := ""
	if t.Upstream.Connection == nil {
		ok = false
		msg += "No connection to upstream server. "
	} else if t.Emitter == nil {
		ok = false
		msg += "No emitter for events. "
	}
	return ok, msg
}

var (
	NAME = "forwarder-tcp-client"
)
