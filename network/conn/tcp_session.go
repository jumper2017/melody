package conn

import (
	"io"
	"net"

	"errors"

	"sync/atomic"

	log "github.com/Sirupsen/logrus"
)

type TcpSession struct {
	conn        net.TCPConn
	sessionId   string      //会话id
	sessionType SessionType //服务端还是客户端
	localAddr   net.Addr    //本端地址
	remoteAddr  net.Addr    //对端地址

	recvChan chan interface{} //接收数据的放置通道

	CbCleanSession func([]string) error //关闭链接的回调函数，　由上层sessionManager进行操作

	//todo: 根据业务需求进行设置，　默认在创建时设置为自身sessionId
	needToCleanSessionsId []string

	closeTag int32 //关闭标记，　防止主动关闭时调用两次　Close
}

func NewTcpSession(sessionId string, sessionType SessionType, conn net.TCPConn, recvChan chan interface{}) (*TcpSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	tcpSession := &TcpSession{
		sessionId:             sessionId,
		sessionType:           sessionType,
		conn:                  conn,
		localAddr:             laddr,
		remoteAddr:            raddr,
		recvChan:              recvChan,
		closeTag:              0,
		needToCleanSessionsId: []string{sessionId},
	}

	return tcpSession, nil
}

func (self *TcpSession) SetCbCleanSession(CbCleanSession func([]string) error) error {
	if CbCleanSession != nil {
		log.Errorf("set cb clean session failed, invalid param.")
		return errors.New("set cb clean session failed, invalid param.")
	}

	self.CbCleanSession = CbCleanSession
	return nil

}

func (self *TcpSession) Start() error {

	for {

		data := make([]byte, STREAM_MSG_HEAD_LENGTH)
		n, err := io.ReadFull(&self.conn, data)

		if atomic.LoadInt32(&self.closeTag) == 1 {
			log.Debugf("close tag is 1, return func.")
			return nil
		}

		if err == io.EOF {
			log.Debugf("recv eof from remote.")
			self.CbCleanSession(self.needToCleanSessionsId)
			return nil
		}

		if err != nil || n != STREAM_MSG_HEAD_LENGTH {
			//todo: 错误处理
			log.Errorf("read data failed, err: %v", err)
			return err
		}

		//转换为长度

		var length int
		data = make([]byte, length)
		n, err = io.ReadFull(&self.conn, data)
		if err != nil || n != STREAM_MSG_HEAD_LENGTH {
			//todo: 错误处理
			log.Errorf("read data failed, err: %v", err)
			return err
		}

		//将数据转到上层
		self.recvChan <- data

	}
	return nil
}

func (self *TcpSession) Send(data []byte) error {

	for {
		n, err := self.conn.Write(data)
		if err != nil {
			log.Errorf("send data failed, err: %v", err)
			return err
		}

		if n == len(data) {
			break
		}

		data = data[n:]
	}

	return nil
}

func (self *TcpSession) Close() error {

	if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
		err := self.conn.Close()
		return err
	}
	return nil
}

func (self *TcpSession) GetSessionId() string {
	return self.sessionId
}
