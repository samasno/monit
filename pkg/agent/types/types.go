package types

import (
	"crypto/tls"
	"net"
)

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
	StdOut(msg string) error
	StdErr(msg string) error
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
