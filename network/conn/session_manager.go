package conn

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"net"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/jumper2017/melody/network/conn/proto"
	"github.com/jumper2017/melody/network/interf"
)

type SessionManager struct {
	// 使用字符串来作为 sessionType, 由三个部分组成 tpc:localip:localport:remoteip:remoteport
	// 即  对应value 的第一个元素

	//考虑，将tcp localip localport remoteip remoteport 作为Session实现的数据成员
	//map key 使用名称比如　"gate:mahjong" 即　"module_name_src:module_name_dst", 即便两个module之间具有多个链接，也可以
	//采用不同名称来进行区别，　名称中不要带协议，　否则就不透明，　没有封装的意义了
	// 这里考虑同一类session可以存在多个数量(连接可能不会具有完全相同的两端ip/port，但是流是可以具有相同的两端ip/port的)，
	// "gate:mahjong/0"表示sessionID, 即链接gate和mahjong的链接的index 0 的session
	sessions     map[string][]interf.Session
	sessionsLock sync.Mutex

	recvChan chan []byte
}

// sessionID 为 "module_name_src:module_name_dst/1" 表示删除 module_name_src:module_name_dst 的第二个元素，
// "module_name_src:module_name_dst/*" 表示删除 module_name_src:module_name_dst 的所有元素
// 同 DelSession 参数含义， 返回可能是一个或多个Session
//CloseSession(sessionID string) error //同 DelSession 参数含义

//session 需要在　client/server　两端进行接口的统一，　因为某些逻辑在 client/server两端是存在差异的，　因此在session内部应该知道自己
//是服务端还是客户端
func NewSession(sessionType SessionType, conn interface{}, recvChan chan []byte) (interf.Session, error) {

	var session interf.Session
	var err error
	switch t := conn.(type) {
	case net.TCPConn:
		session, err = NewTcpSession(sessionType, t, recvChan)
		break

	case websocket.Conn:
		session, err = NewWsSession(sessionType, t, recvChan)
		break

	case proto.GrpcBid_CommServer:
		session, err = NewGBSSession(sessionType, t, recvChan)
		break

	case net.UDPConn:
		session, err = NewUdpSession(sessionType, t, recvChan)
		break

	default:
		session, err = nil, nil
		break
	}

	return session, err
}

func (self *SessionManager) AddSession(sessionId string, s interf.Session) error {

	if s == nil {
		return errors.New("add session failed | invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	if slice, ok := self.sessions[sessionId]; ok {
		slice = append(slice, s)
	} else {
		self.sessions[sessionId] = []interf.Session{s}
	}

	s.SetCbCleanSession(self.CleanSession)
	newSessionId := fmt.Sprintf("%s/%d", sessionId, len(self.sessions[sessionId])-1)
	s.SetSessionId(newSessionId)

	return nil
}

func (self *SessionManager) DelSession(sessionID string) error {

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	return self.delSessionWithoutLock(sessionID)
}

func (self *SessionManager) delSessionWithoutLock(sessionID string) error {

	sname := strings.Split(sessionID, "/")
	if len(sname) != 2 {
		return errors.New("del session failed | invalid param.")
	}

	if sname[1] == "*" {
		delete(self.sessions, sname[0])
	} else {
		tmpslice := self.sessions[sname[0]]
		index, err := strconv.Atoi(sname[1])
		if err != nil {
			return err
		}
		tmpslice = append(tmpslice[:index], tmpslice[index+1:]...)
		self.sessions[sname[0]] = tmpslice
	}

	return nil
}

func (self *SessionManager) GetSession(sessionID string) ([]interf.Session, error) {

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	return self.getSessionWithoutLock(sessionID)
}

func (self *SessionManager) getSessionWithoutLock(sessionID string) ([]interf.Session, error) {

	sname := strings.Split(sessionID, "/")
	if len(sname) != 2 {
		return nil, errors.New("get session failed | invalid param.")
	}

	if sname[1] == "*" {
		return self.sessions[sname[0]], nil
	} else {
		index, err := strconv.Atoi(sname[1])
		if err != nil {
			return nil, err
		}
		return []interf.Session{self.sessions[sname[0]][index]}, nil
	}

	return nil, errors.New("get session failed | invalid param.")
}

func (self *SessionManager) CloseSession(sessionID string) error {

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	return self.closeSessionWithoutLock(sessionID, true)
}

func (self *SessionManager) closeSessionWithoutLock(sessionID string, method bool) error {

	sname := strings.Split(sessionID, "/")
	if len(sname) != 2 {
		return errors.New("get session failed | invalid param.")
	}

	if sname[1] == "*" {
		for _, v := range self.sessions[sname[0]] {
			if !method {
				v.ClosePassive()
			} else {
				v.CloseInitiative()
			}
		}
		return nil
	} else {
		index, err := strconv.Atoi(sname[1])
		if err != nil {
			return err
		}

		if !method {
			self.sessions[sname[0]][index].ClosePassive()
		} else {
			self.sessions[sname[0]][index].CloseInitiative()
		}
		return nil
	}

	return errors.New("close session failed | invalid param.")
}

//需要考虑session关闭的多种需求.
//1. 一个链接关闭，但是还有其他链接在引用该recvChan，　此时无需关闭 recvChan，　只需要关闭该链接，　（比如退出房间的情况）
//2. 一个链接关闭，需要关闭sessionManager下的所有其他链接，（比如client链接关闭了，此时需要关闭gate-room链接以及gate-xxx链接）
//3. 一个链接关闭，需要关闭sessionManager下的部分链接，因此需要指定需要关闭的链接的sessionId
//上述3种需求可以统一为一个接口
func (self *SessionManager) CleanSession(sessionIds []string, method bool) error {
	//sessionIds 可能取值：
	// 1. 一个需要关闭的sessionId, 即底层调用该回调函数的session, 如 gate:mahjong/0
	// 2. "all" 表示关闭所有链接
	// 3. "sessionId1, sessionId2" 表示需要关闭的链接, 包含底层调用该回调函数的session

	if sessionIds == nil {
		log.Errorf("clean session failed, invalid param.")
		return errors.New("clean session failed, invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	if len(sessionIds) == 1 && sessionIds[0] == "all" {
		for k := range self.sessions {
			self.closeSessionWithoutLock(k+"/*", method)
			self.delSessionWithoutLock(k + "/*")
		}
	} else {
		for _, v := range sessionIds {
			self.closeSessionWithoutLock(v, method)
			self.delSessionWithoutLock(v)
		}
	}

	return nil
}

//用于关闭recvChan，　需要小心，　调用该函数需要保证之后没有session会再向 recvChan 中发送数据
//在何处调用该函数应该取决于具体的业务需求
//todo: 这里需要根据具体的需求来进行考虑，　先放一放
func (self *SessionManager) CloseRecvChan() error {

	close(self.recvChan)
	return nil
}
