package types

import (
	"crypto/tls"
	"net"
)

type Controller interface { // manages the forwarder to upstream and log runners
	Init() error
	Run() error
	Shutdown() error
	Status() (Status, error)
}

type Forwarder interface {
	Connect() error
	Close() error
	Push([]byte) error
	Status() (Status, error)
}

type LogTail interface {
	Open() error
	Close() error
	Update() error
	Status() (Status, error)
}

type Logger interface {
	StdOut() error
	StdErr() error
	Close() error
	Status() (Status, error)
}

type Status interface {
	Message() string
	Ok() bool
}

type Upstream struct {
	Connection net.Conn
	URL        string
	TlsConfig  *tls.Config
}
