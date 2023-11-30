package i

import (
	t "github.com/samasno/monit/pkg/agent/types"
)

type Controller interface { // manages the forwarder to upstream and log runners
	Init(input t.ControllerInitInput) error
	Run()
	Shutdown()
	Status() t.ControllerStatus
}

type Forwarder interface {
	Connect()
	Close()
	Push()
	Status()
}

type LogTail interface {
	Open()
	Close()
	Update()
	Status()
}

type Logger interface {
	Log([]byte)
}
