package conn

import (
	"errors"
	"net"
)

type TcpSession struct {
	conn       net.TCPConn
	sessionId  uint32   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道
}

func (self *TcpSession) Send(data []byte) error {

	return nil
}

//func (self *TcpSession) Recv() ([]byte, error) {
//
//	return nil, nil
//}

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

func (self *TcpSession) SetRecvChan(chan interface{}) error {

}

func (self *TcpSession) Start() error {

	return nil
}
