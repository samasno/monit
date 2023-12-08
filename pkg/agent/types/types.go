package types

import (
	"crypto/tls"
	"net"
)

type Forwarder interface {
	Connect() error
	Close() error
	Push([]byte) error
}

type ForwarderClient interface {
	Connect() error
	Disconnect() error
	Push([]byte) error
}

type ForwarderListener interface {
	Open() error
	Close() error
	Listen() ([]byte, error)
}

type Emitter interface {
	Emit(Event) error
}

type LogTail interface {
	Open() error
	Close() error
	Update() error
}

type Logger interface {
	StdOut(msg string) error
	StdErr(msg string) error
	Close() error
}

type Status interface {
	Message() string
}

type Upstream struct {
	Connection net.Conn
	Url        string
	Port       int
	TlsConfig  *tls.Config
}

type Downstream struct {
	Connection net.Conn
	Url        string
	Port       int
	TlsConfig  *tls.Config
}

type Event struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload,omitempty"`
}

type Payload struct {
	Source  string `json:"source,omitempty"`
	Message string `json:"message,omitempty"`
	Level   int    `json:"level,omitempty"`
}
