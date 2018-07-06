package conn

import (
	"errors"
	"net"
)

type UdpSession struct {
	conn       net.UDPConn
	sessionId  uint32   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道

}

func (self *UdpSession) Send(data []byte) error {

	return nil
}

//func (self *UdpSession) Recv() ([]byte, error) {
//
//	return nil, nil
//}

func (self *UdpSession) Close() error {

	return self.conn.Close()
}

func (self *UdpSession) SetRawConn(conn interface{}) error {

	if c, ok := conn.(net.UDPConn); ok {
		self.conn = c
		return nil
	}

	return errors.New("set raw conn in UdpSession failed | invalid param.")
}

func (self *UdpSession) SetRecvChan(chan interface{}) error {

}

func (self *UdpSession) Start() error {

	return nil
}
