package emitters

import (
	"encoding/json"
	"net"

	"github.com/samasno/monit/pkg/agent/types"
)

type SocketEmitter struct {
	conn  *net.UnixConn
	Raddr string
	Laddr string
}

func (s *SocketEmitter) Emit(event types.Event) error {
	if s.conn == nil {
		err := s.connect()
		if err != nil {
			s.conn = nil
			return err
		}
	}
	e, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(e)
	if err != nil {
		return err
	}
	return nil
}

func (s *SocketEmitter) connect() error {
	if s.conn != nil {
		return nil
	}
	laddr, err := net.ResolveUnixAddr("unixgram", s.Laddr)
	if err != nil {
		return err
	}
	raddr, err := net.ResolveUnixAddr("unixgram", s.Raddr)
	if err != nil {
		return err
	}
	cn, err := net.DialUnix("unixgram", laddr, raddr)
	if err != nil {
		return err
	}
	cn.SetWriteBuffer(65536)
	s.conn = cn
	return nil
}
