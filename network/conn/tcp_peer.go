package conn

import (
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/jumper2017/melody/network/interf"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////
type TcpPeerAcceptor struct {
	funcAddSession func(sessionName string, s interf.Session)
	listener       net.Listener
}

//接收到一个链接请求之后， 会创建一个session，
// 调用f以便将session 传入到对应的agent中
func (self *TcpPeerAcceptor) RegisterGenerateSession(f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.funcAddSession = f
	return
}

func (self *TcpPeerAcceptor) Start(sessionName string, listenAddr string, recvChan chan []byte) {

	var err error
	self.listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		panic("start tcp acceptor failed.")
	}

	for {
		conn, err := self.listener.Accept()

		if err != nil {
			fmt.Println("err:", err)
			logrus.Errorf("accept failed. err: %v", err)
			return
		}

		tcpSession, err := NewSession(ServerSession, conn.(*net.TCPConn), recvChan)
		if err != nil {
			logrus.Errorf("new tcp session failed. err: %v", err)
			return
		}

		go self.funcAddSession(sessionName, tcpSession)
	}

}

func (self *TcpPeerAcceptor) Stop() bool {
	self.listener.Close()
	return true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
type TcpPeerConnector struct {
	// 通过主动链接获得session之后， 存放到agent 的 session manager 中
	funcAddSession func(sessionName string, s interf.Session)
}

func (self *TcpPeerConnector) RegisterGenerateSession(f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.funcAddSession = f
	return
}

func (self *TcpPeerConnector) Start(sessionName string, connAddr string, recvChan chan []byte) {

	conn, err := net.Dial("tcp", connAddr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	tcpSession, err := NewSession(ClientSession, conn.(*net.TCPConn), recvChan)
	if err != nil {
		logrus.Errorf("new tcp session failed. err: %v", err)
		return
	}

	go self.funcAddSession(sessionName, tcpSession)

	return
}
