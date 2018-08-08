package conn

import (
	"fmt"

	"net/http"

	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/jumper2017/melody/network/interf"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}
)

type WsPeerAcceptor struct {
	funcAddSession func(sessionName string, s interf.Session)
}

//接收到一个链接请求之后， 会创建一个session，
// 调用f以便将session 传入到对应的agent中
func (self *WsPeerAcceptor) RegisterGenerateSession(f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.funcAddSession = f
	return
}

func (self *WsPeerAcceptor) Start(sessionName string, listenAddr string, recvChan chan []byte) {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Errorf("upgrade failed.")
			return
		}

		wsSession, err := NewSession(ServerSession, conn, recvChan)
		if err != nil {
			logrus.Errorf("new ws session failed. err: %v", err)
			return
		}

		self.funcAddSession(sessionName, wsSession)

		return
	})

	http.ListenAndServe(listenAddr, nil)
	return
}

func (self *WsPeerAcceptor) Stop() bool {
	return true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
type WsPeerConnector struct {
	// 通过主动链接获得session之后， 存放到agent 的 session manager 中
	funcAddSession func(sessionName string, s interf.Session)
	connAddr       string
}

func (self *WsPeerConnector) InitConn(connAddr string, f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.connAddr = connAddr
	self.funcAddSession = f
	return
}

func (self *WsPeerConnector) Start(sessionName string, recvChan chan []byte) {

	u := url.URL{Scheme: "ws", Host: self.connAddr, Path: "/ws"}
	var dia websocket.Dialer
	conn, _, err := dia.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("error in dial, err:", err)
		return
	}

	wsSession, err := NewSession(ClientSession, conn, recvChan)
	if err != nil {
		logrus.Errorf("new ws session failed. err: %v", err)
		return
	}

	go self.funcAddSession(sessionName, wsSession)
	return
}
