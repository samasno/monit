package forwarder

import (
	"fmt"

	types "github.com/samasno/monit/pkg/agent/types"
)

type Forwarder struct {
	Name           string
	UpstreamClient types.ForwarderClient
	EventListener  types.ForwarderListener
	Ok             bool
}

func (f *Forwarder) Connect() error {
	err := f.UpstreamClient.Connect()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to connect to upstream %s\n", f.Name, err.Error())
	}
	err = f.EventListener.Open()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to connect %s\n", f.Name, err.Error())
	}
	return nil
}

func (f *Forwarder) Close() error {
	err := f.UpstreamClient.Disconnect()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to close connection to upstream %s\n", f.Name, err.Error())
	}
	err = f.EventListener.Close()
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to close listener", f.Name)
	}
	return nil
}

func (f *Forwarder) Push(msg []byte) error {
	err := f.UpstreamClient.Push(msg)
	if err != nil {
		return fmt.Errorf("Forwarder %s: Failed to push to upstream %s\n", f.Name, err.Error())
	}
	return nil
}

func (f *Forwarder) Status() (types.Status, error) {
	fs := ForwarderStatus{}
	if f.UpstreamClient == nil {
		fs.IsOk = false
		fs.MessageText = "Missing upstream client. "
	}

	if f.EventListener == nil {
		fs.IsOk = false
		fs.MessageText += "No listening client"
	}

	if fs.MessageText == "" {
		fs.MessageText = "Looks ok"
		fs.IsOk = true
	}
	return fs, nil
}
