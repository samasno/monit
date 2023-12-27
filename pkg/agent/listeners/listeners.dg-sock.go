package listeners

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	types "github.com/samasno/monit/pkg/agent/types"
	"github.com/samasno/monit/pkg/agent/vars"
)

var defaultUDGSName = "unix-datagram-socket-listener"

type UnixDatagramSocketListener struct {
	Name       string
	Downstream *types.Downstream
	Logger     types.Emitter
	shutdown   *sync.WaitGroup
}

func (l *UnixDatagramSocketListener) Open(shutdown *sync.WaitGroup) error {
	if l.shutdown == nil {
		l.shutdown = shutdown
		l.shutdown.Add(1)
	}
	if l.Downstream.Connection != nil {
		l.log(vars.NOTICE, "Already listening at "+l.Downstream.Url)
		return nil
	}
	addr, err := net.ResolveUnixAddr("unixgram", l.Downstream.Url)
	ln, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		msg := fmt.Sprintf("Failed to open unix datagram socket at " + l.Downstream.Url)
		l.log(vars.CRITICAL, msg)
		return fmt.Errorf("%s: %s\n", msg, err.Error())
	}
	ln.SetWriteBuffer(65536)
	ln.SetReadBuffer(65536)
	l.log(vars.INFO, "Set socket write buffer to 65536")
	l.log(vars.INFO, "Set socket read buffer to 65536")
	l.Downstream.Connection = ln
	l.log(vars.INFO, "Listening on unix datagram socket at "+l.Downstream.Url)
	return nil
}

func (l *UnixDatagramSocketListener) Close() error {
	if l.Downstream.Connection == nil {
		l.log(vars.NOTICE, "Socket at "+l.Downstream.Url+" is already closed.")
		return nil
	}
	l.log(vars.INFO, "Closing unix datagram socket at "+l.Downstream.Url)
	err := l.Downstream.Connection.Close()
	if err != nil {
		msg := fmt.Sprintf("Failed to close unix datagram socket at %s: %s", l.Downstream.Url, err.Error())
		l.log(vars.ERROR, msg)
		return fmt.Errorf(msg + "\n")
	}
	l.Downstream.Connection = nil
	os.Remove(l.Downstream.Url)
	if l.shutdown != nil {
		l.shutdown.Done()
		l.shutdown = nil
	}
	l.log(vars.INFO, fmt.Sprintf("Closed unix datagram socket at %s", l.Downstream.Url))
	return nil
}

func (l *UnixDatagramSocketListener) Listen(out chan []byte, closer chan bool, shutdown *sync.WaitGroup) error {
	if l.Downstream.Connection == nil {
		err := l.Open(shutdown)
		if err != nil {
			msg := "Failed to open downstream listener"
			l.log(vars.ERROR, msg)
			return err
		}
	}
	go func(out chan []byte, closer chan bool) {
		running := true
		go func(out chan []byte) {
			defer func() {
				if r := recover(); r != nil {
					l.log(vars.ERROR, "Recovered from panic in Listen")
				}
			}()
			for {
				b := make([]byte, 4096)
				if !running {
					break
				}
				n, err := l.Downstream.Connection.Read(b)
				if err != nil {
					if _, ok := err.(*net.OpError); ok && running {
						l.log(vars.CRITICAL, "Socket has closed unexpectedly")
						l.log(vars.NOTICE, "Restarting socket in 5 seconds")
						time.Sleep(5 * time.Second)
						l.Listen(out, closer, shutdown)
						break
					} else {
						if running {
							l.log(vars.ERROR, "Failed to read packet")
						} else {
							l.log(vars.NOTICE, "No longer accepting packets")
						}
					}
				}
				if n > 0 {
					select {
					case out <- b[:n]:
						l.log(vars.INFO, fmt.Sprintf("Sending %d bytes to forward", n))
					default:
						l.log(vars.NOTICE, "Looks like out buffer is full, message dropped")
					}
				}
			}
		}(out)
		for {
			done := <-closer
			if done {
				running = false
				l.log(vars.NOTICE, "Received close signal")
				err := l.Close()
				if err != nil {
					l.log(vars.ERROR, "Failed to close unix datagram socket")
				}
				if err != nil {
					println("not running is true :" + err.Error())
				}
				closer <- true
				break
			}
		}
	}(out, closer)
	return nil
}

func (l *UnixDatagramSocketListener) log(level int, message string) error {
	if l.Logger == nil {
		return fmt.Errorf("No emitter to send logs\n")
	}
	var name string
	if l.Name == "" {
		name = defaultUDGSName
	} else {
		name = l.Name
	}
	payload := types.Payload{
		Source:  name,
		Message: message,
		Level:   level,
	}
	event := types.Event{
		Type:    vars.LISTENER_CLIENT_LOG,
		Payload: payload,
	}
	err := l.Logger.Emit(event)
	if err != nil {
		msg := fmt.Sprintf("Failed to emit message")
		return fmt.Errorf(msg + "\n")
	}
	return nil
}
