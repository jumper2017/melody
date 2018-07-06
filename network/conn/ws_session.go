package conn

import (
	"errors"
	"net"

	"github.com/gorilla/websocket"
)

type WsSession struct {
	conn       websocket.Conn
	sessionId  uint32   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道

}

func (self *WsSession) Send(data []byte) error {

	return nil
}

//func (self *WsSession) Recv() ([]byte, error) {
//
//	return nil, nil
//}

func (self *WsSession) Close() error {

	return self.conn.Close()
}

func (self *WsSession) SetRawConn(conn interface{}) error {

	if c, ok := conn.(websocket.Conn); ok {
		self.conn = c
		return nil
	}

	return errors.New("set raw conn in WsSession failed | invalid param.")
}

func (self *WsSession) SetRecvChan(chan interface{}) error {

}

func (self *WsSession) Start() error {

	return nil
}
