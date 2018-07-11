package conn

import (
	"net"

	"github.com/gorilla/websocket"
	"fmt"
)

type WsSession struct {
	conn       websocket.Conn
	sessionId  string   //会话id
	localAddr  net.Addr //本端地址
	remoteAddr net.Addr //对端地址

	recvChan chan interface{} //接收数据的放置通道

}

func NewWsSession(conn websocket.Conn, recvChan chan interface{}) (*WsSession, error){

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	wsSession := &WsSession{
		conn: conn,
		localAddr: laddr,
		remoteAddr: raddr,
		recvChan: recvChan,
		sessionId: fmt.Sprintf("%s:%s:%s", "tcp", laddr.String(), raddr.String()),
	}

	return wsSession, nil
}

func (self *WsSession) Start() error {

	return nil
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


