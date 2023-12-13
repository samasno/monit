package forwarder

import (
	"fmt"
	"sync"

	types "github.com/samasno/monit/pkg/agent/types"
)

type Forwarder struct {
	UpstreamClient     types.ForwarderClient
	DownstreamListener types.ForwarderListener
	listenerCloser     chan bool
	listenerOutput     chan []byte
	shutdown           *sync.WaitGroup
}

func (f *Forwarder) Connect() error {
	f.listenerCloser = make(chan bool)
	f.listenerOutput = make(chan []byte, 100)
	f.shutdown = &sync.WaitGroup{}
	err := f.UpstreamClient.Connect(f.shutdown)
	if err != nil {
		return fmt.Errorf("Forwarder: Failed to connect to upstream %s\n", err.Error())
	}

	err = f.DownstreamListener.Listen(f.listenerOutput, f.listenerCloser, f.shutdown)
	if err != nil {
		return fmt.Errorf("Forwarder: Failed to connect %s\n", err.Error())
	}
	return nil
}

func (f *Forwarder) Close() error {
	f.listenerCloser <- true
	close(f.listenerOutput)
	f.UpstreamClient.Disconnect()
	f.shutdown.Wait()
	return nil
}

func (f *Forwarder) Push(msg []byte) error {
	err := f.UpstreamClient.Push(msg)
	if err != nil {
		return fmt.Errorf("Forwarder: Failed to push to upstream %s\n", err.Error())
	}
	return nil
}

func (f *Forwarder) Run() error {
	f.listenerCloser = make(chan bool)
	f.listenerOutput = make(chan []byte, 100)
	err := f.Connect()
	if err != nil {
		return err
	}
	err = f.DownstreamListener.Listen(f.listenerOutput, f.listenerCloser, f.shutdown)
	go func() {
		var running = true
		for {
			select {
			case output, ok := <-f.listenerOutput:
				if !ok {
					running = false
				} else {
					f.UpstreamClient.Push(output)
				}
			}
			if !running {
				break
			}
		}
	}()
	return nil
}

func (f *Forwarder) Status() (types.Status, error) {
	fs := ForwarderStatus{}
	if f.UpstreamClient == nil {
		fs.IsOk = false
		fs.MessageText = "Missing upstream client. "
	}

	if f.DownstreamListener == nil {
		fs.IsOk = false
		fs.MessageText += "No listening client"
	}

	if fs.MessageText == "" {
		fs.MessageText = "Looks ok"
		fs.IsOk = true
	}
	return fs, nil
}
