package load_balancer

import (
	"errors"
	"net"
)

type connListener struct {
	conn  net.Conn
	close chan struct{}
}

func (c connListener) Accept() (conn net.Conn, err error) {
	if c.conn != nil {
		conn = c.conn
		c.conn = nil
		return
	}

	<-c.close
	err = errors.New("connection listener is closed")
	return
}

func (c connListener) Close() error {
	close(c.close)
	//_ = c.conn.Close()
	return nil
}

func (c connListener) Addr() net.Addr {
	return c.conn.LocalAddr()
}
