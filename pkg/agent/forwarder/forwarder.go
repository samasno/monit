package forwarder

import (
	"fmt"

	types "github.com/samasno/monit/pkg/agent/types"
)

type Forwarder struct {
	Name           string
	upstreamClient types.ForwarderClient
	eventListener  types.ForwarderListener
	Ok             bool
}

func (f *Forwarder) Connect() error {
	err := f.upstreamClient.Connect()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to connect to upstream %s\n", f.Name, err.Error())
	}
	err = f.eventListener.Open()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to connect %s\n", f.Name, err.Error())
	}
	return nil
}

func (f *Forwarder) Close() error {
	err := f.upstreamClient.Disconnect()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to close connection to upstream %s\n", f.Name, err.Error())
	}
	err = f.eventListener.Close()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to close listener", f.Name)
	}
	return nil
}

func (f *Forwarder) Push(msg []byte) error {
	err := f.upstreamClient.Push(msg)
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to push to upstream %s\n", f.Name, err.Error())
	}
	return nil
}

func (f *Forwarder) Status() (types.Status, error) {
	fs := ForwarderStatus{}
	if f.upstreamClient == nil {
		fs.IsOk = false
		fs.MessageText = "Missing upstream client. "
	}

	if f.eventListener == nil {
		fs.IsOk = false
		fs.MessageText += "No listening client"
	}

	if fs.MessageText == "" {
		fs.MessageText = "Looks ok"
		fs.IsOk = true
	}
	return fs, nil
}
