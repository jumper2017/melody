package conn

import (
	"net"
	"fmt"
)

type TcpSession struct {
	conn       net.TCPConn
	sessionId  string   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道
}

func NewTcpSession(conn net.TCPConn, recvChan chan interface{}) (*TcpSession, error){

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	tcpSession := &TcpSession{
		conn: conn,
		localAddr: laddr,
		remoteAddr: raddr,
		recvChan: recvChan,
		sessionId: fmt.Sprintf("%s:%s:%s", "tcp", laddr.String(), raddr.String()),
	}

	return tcpSession, nil
}

func (self *TcpSession) Start() error {

	return nil
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


