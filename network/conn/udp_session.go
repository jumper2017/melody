package conn

import (
	"net"
	"fmt"
)

type UdpSession struct {
	conn       net.UDPConn
	sessionId  string   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道

}

func NewUdpSession(conn net.UDPConn, recvChan chan interface{}) (*UdpSession, error){

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	udpSession := &UdpSession{
		conn: conn,
		localAddr: laddr,
		remoteAddr: raddr,
		recvChan: recvChan,
		sessionId: fmt.Sprintf("%s:%s:%s", "tcp", laddr.String(), raddr.String()),
	}

	return udpSession, nil
}

func (self *UdpSession) Start() error {

	return nil
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


