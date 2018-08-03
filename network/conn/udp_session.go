package conn

import (
	"errors"
	"net"

	log "github.com/Sirupsen/logrus"
)

type UdpSession struct {
	conn        net.UDPConn
	sessionId   string      //会话id
	sessionType SessionType //服务端还是客户端
	localAddr   net.Addr    //本端地址
	remoteAddr  net.Addr    //对端地址

	recvChan chan []byte //接收数据的放置通道

	CbCleanSession func([]string) error //关闭链接的回调函数，　由上层sessionManager进行操作
}

func NewUdpSession(sessionType SessionType, conn net.UDPConn, recvChan chan []byte) (*UdpSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	udpSession := &UdpSession{
		sessionType: sessionType,
		conn:        conn,
		localAddr:   laddr,
		remoteAddr:  raddr,
		recvChan:    recvChan,
	}

	return udpSession, nil
}

func (self *UdpSession) SetCbCleanSession(CbCleanSession func([]string) error) error {
	if CbCleanSession != nil {
		log.Errorf("set cb clean session failed, invalid param.")
		return errors.New("set cb clean session failed, invalid param.")
	}

	self.CbCleanSession = CbCleanSession
	return nil

}

func (self *UdpSession) Start() error {

	for {
		self.conn.Read()
	}
	return nil
}

func (self *UdpSession) Send(data []byte) error {

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

func (self *UdpSession) Close() error {

	return self.conn.Close()
}

func (self *UdpSession) GetSessionId() string {
	return self.sessionId
}

func (self *UdpSession) SetSessionId(sessionId string) {
	self.sessionId = sessionId
	self.needToCleanSessionsId = []string{sessionId}
	return
}
