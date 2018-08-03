package conn

import (
	"errors"

	"sync/atomic"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type WsSession struct {
	conn     websocket.Conn
	recvChan chan []byte //接收数据的放置通道
	BaseSession
}

func NewWsSession(sessionType SessionType, conn websocket.Conn, recvChan chan []byte) (*WsSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	wsSession := &WsSession{

		BaseSession: BaseSession{
			sessionType: sessionType,
			localAddr:   laddr,
			remoteAddr:  raddr,
			closeTag:    0,
		},
		conn:     conn,
		recvChan: recvChan,
	}

	return wsSession, nil
}

func (self *WsSession) Start() error {

	for {

		msgType, data, err := self.conn.ReadMessage()

		if atomic.LoadInt32(&self.closeTag) == 1 {
			log.Debugf("close tag is 1, return func.")
			return nil
		}

		switch {

		case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
			log.Debugf("recv eof from remote.")
			self.CbCleanSession(self.needToCleanSessionsId, false)
			return nil

		case websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
			log.Debugf("recv err eof from remote, %v", err)
			self.CbCleanSession(self.needToCleanSessionsId, false)
			return errors.New("recv err eof from remote")

		case err != nil:
			log.Errorf("read msg failed, err: %v", err)
			self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, "server close connection."))
			self.CbCleanSession(self.needToCleanSessionsId, true)
			return errors.New("read msg failed")

			//todo: 根据业务需求进行修改，　在对端不是使用该golang库的情况下，不需要进行如下判断
		case msgType != websocket.BinaryMessage:
			log.Errorf("invalid msg type: %v", msgType)
			self.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "server close connection."))
			self.CbCleanSession(self.needToCleanSessionsId, true)
			return errors.New("invalid msg type")
		}

		self.recvChan <- data

	}

	return nil
}

func (self *WsSession) Send(data []byte) error {

	return self.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (self *WsSession) ClosePassive() error {

	if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
		err := self.conn.Close()
		if err != nil {
			log.Errorf("close failed, err: %v", err)
			return err
		}
	}
	return nil

}

func (self *WsSession) CloseInitiative() error {
	return self.ClosePassive()
}
