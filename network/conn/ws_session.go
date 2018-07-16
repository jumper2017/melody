package conn

import (
	"errors"
	"net"

	"sync/atomic"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type WsSession struct {
	conn        websocket.Conn
	sessionId   string      //会话id
	sessionType SessionType //服务端还是客户端
	localAddr   net.Addr    //本端地址
	remoteAddr  net.Addr    //对端地址

	recvChan chan interface{} //接收数据的放置通道

	CbCleanSession func([]string) error //关闭链接的回调函数，　由上层sessionManager进行操作

	closeTag int32 //关闭动作标记

	//todo: 根据业务需求进行设置，　默认在创建时设置为自身sessionId
	needToCleanSessionsId []string
}

func NewWsSession(sessionId string, sessionType SessionType, conn websocket.Conn, recvChan chan interface{}) (*WsSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	wsSession := &WsSession{
		sessionId:             sessionId,
		sessionType:           sessionType,
		conn:                  conn,
		localAddr:             laddr,
		remoteAddr:            raddr,
		recvChan:              recvChan,
		closeTag:              0,
		needToCleanSessionsId: []string{sessionId},
	}

	return wsSession, nil
}

func (self *WsSession) SetCbCleanSession(CbCleanSession func([]string) error) error {
	if CbCleanSession != nil {
		log.Errorf("set cb clean session failed, invalid param.")
		return errors.New("set cb clean session failed, invalid param.")
	}

	self.CbCleanSession = CbCleanSession
	return nil

}

func (self *WsSession) Start() error {

	for {

		msgType, data, err := self.conn.ReadMessage()

		if atomic.LoadInt32(&self.closeTag) == 1 {
			log.Debugf("close tag is 1, return func.")
			return nil
		}

		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Debugf("recv eof from remote.")
			self.CbCleanSession(self.needToCleanSessionsId)
			return nil
		}

		if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			log.Debugf("recv err eof from remote, %v", err)
			self.CbCleanSession(self.needToCleanSessionsId)
			return errors.New("recv err eof from remote")
		}

		if err != nil {
			log.Errorf("read msg failed, err: %v", err)
			self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, "server close connection."))
			self.CbCleanSession(self.needToCleanSessionsId)
			return errors.New("read msg failed")
		}

		//todo: 根据业务需求进行修改，　在对端不是使用该golang库的情况下，不需要进行如下判断
		if msgType != websocket.BinaryMessage {
			log.Errorf("invalid msg type: %v", msgType)
			self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "server close connection."))
			self.CbCleanSession(self.needToCleanSessionsId)
			return errors.New("invalid msg type")
		}

		self.recvChan <- data

	}

	return nil
}

func (self *WsSession) Send(data []byte) error {

	//todo: 根据业务需求进行修改
	self.conn.WriteMessage(websocket.BinaryMessage, data)
	return nil
}

func (self *WsSession) Close() error {

	if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
		err := self.conn.Close()
		if err != nil {
			log.Errorf("close failed, err: %v", err)
			return err
		}
	}
	return nil

}

func (self *WsSession) GetSessionId() string {
	return self.sessionId
}
