package conn

import (
	"net"

	"fmt"

	"github.com/gorilla/websocket"
)

type WsSession struct {
	conn        websocket.Conn
	sessionId   string      //会话id
	sessionType SessionType //服务端还是客户端
	localAddr   net.Addr    //本端地址
	remoteAddr  net.Addr    //对端地址

	recvChan chan interface{} //接收数据的放置通道

}

func NewWsSession(sessionType SessionType, conn websocket.Conn, recvChan chan interface{}) (*WsSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	wsSession := &WsSession{
		conn:        conn,
		localAddr:   laddr,
		remoteAddr:  raddr,
		recvChan:    recvChan,
		sessionType: sessionType,
		sessionId:   fmt.Sprintf("%s:%s:%s", "tcp", laddr.String(), raddr.String()),
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

func (self *WsSession) GetSessionId() string {
	return self.sessionId
}
