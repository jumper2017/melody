package conn

import (
	"io"
	"net"

	"sync/atomic"

	"encoding/binary"

	log "github.com/Sirupsen/logrus"
)

type TcpSession struct {
	conn     *net.TCPConn
	recvChan chan []byte //接收数据的放置通道
	BaseSession
}

func NewTcpSession(sessionType SessionType, conn *net.TCPConn, recvChan chan []byte) (*TcpSession, error) {

	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	tcpSession := &TcpSession{

		BaseSession: BaseSession{
			sessionType: sessionType,
			localAddr:   laddr,
			remoteAddr:  raddr,
			closeTag:    0,
		},
		conn:     conn,
		recvChan: recvChan,
	}

	return tcpSession, nil
}

func (self *TcpSession) Start() error {

	for {

		data := make([]byte, STREAM_MSG_HEAD_LENGTH)
		n, err := io.ReadFull(self.conn, data)

		if atomic.LoadInt32(&self.closeTag) == 1 {
			log.Debugf("close tag is 1, return func.")
			return nil
		}

		switch {

		case err == io.EOF:
			log.Debugf("recv eof from remote.")
			self.CbCleanSession(self.needToCleanSessionsId, false)
			return nil

		case err != nil || n != STREAM_MSG_HEAD_LENGTH:
			log.Errorf("read data failed, err: %v", err)
			self.CbCleanSession(self.needToCleanSessionsId, true)
			return err

		}

		//解出长度
		length := binary.LittleEndian.Uint16(data)
		data = make([]byte, length)
		n, err = io.ReadFull(self.conn, data)
		if err != nil || n != int(length) {
			//todo: 错误处理
			log.Errorf("read data failed, err: %v", err)
			self.CbCleanSession(self.needToCleanSessionsId, true)
			return err
		}

		//将数据转到上层
		self.recvChan <- data

	}
	return nil
}

func (self *TcpSession) Send(data []byte) error {

	//添加length
	length := len(data)
	sendData := make([]byte, length+STREAM_MSG_HEAD_LENGTH)
	binary.LittleEndian.PutUint16(sendData, uint16(length))
	copy(sendData[STREAM_MSG_HEAD_LENGTH:], data)

	for {
		n, err := self.conn.Write(sendData)
		if err != nil {
			log.Errorf("send sendData failed, err: %v", err)
			return err
		}

		if n == len(sendData) {
			break
		}

		sendData = sendData[n:]
	}

	return nil
}

func (self *TcpSession) ClosePassive() error {

	if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
		err := self.conn.Close()
		return err
	}
	return nil
}

func (self *TcpSession) CloseInitiative() error {
	return self.ClosePassive()
}
