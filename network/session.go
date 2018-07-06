package network

import (
	"net"

	"errors"
)

type Session interface {
	SetRawConn(conn interface{})

	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

type TcpSession struct {
	conn       net.TCPConn
	sessionId  uint32 //会话id
	localAddr  string //本端地址
	remoteAddr string //对端地址

}

func (self *TcpSession) Send(data []byte) error {

	return nil
}

func (self *TcpSession) Recv() ([]byte, error) {

	return nil, nil
}

func (self *TcpSession) Close() error {

	return self.conn.Close()
}

func (self *TcpSession) SetRawConn(conn interface{}) error {

	if c, ok := conn.(net.TCPConn); ok {
		self.conn = c
		return nil
	}

	return errors.New("set raw conn in TcpSession failed | invalid param.")
}
