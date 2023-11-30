package t

import (
	"crypto/tls"
	"net"
)

type ControllerInitInput struct {
	UpstreamDetails Upstream
	LogPaths        []string
}

type ControllerStatus struct{}

type Upstream struct {
	Connection net.Conn
	URL        string
	TlsConfig  *tls.Config
}
