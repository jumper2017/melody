package conn

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/jumper2017/melody/network/interf"
)

// sessionID 为 "tpc:localip:localport:remoteip:remoteport/1" 表示删除 tpc:localip:localport:remoteip:remoteport 的第二个元素，
// "tpc:localip:localport:remoteip:remoteport/*" 表示删除 tpc:localip:localport:remoteip:remoteport 的所有元素
// 同 DelSession 参数含义， 返回可能是一个或多个Session
//CloseSession(sessionID string) error //同 DelSession 参数含义
type SessionManager struct {
	// 使用字符串来作为 sessionType, 由三个部分组成 tpc:localip:localport:remoteip:remoteport
	// 这里考虑同一类session可以存在多个数量， "tpc:localip:localport:remoteip:remoteport/0"表示sessionID, 即
	// 链接到gate的tcp协议类型的index 0 的session, 即  对应value 的第一个元素
	sessions     map[string][]interf.Session
	sessionsLock sync.Mutex
}

func (self *SessionManager) AddSession(sessionType string, s interf.Session) error {

	if s == nil {
		return errors.New("add session failed | invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	if slice, ok := self.sessions[sessionType]; ok {
		slice = append(slice, s)
	} else {
		self.sessions[sessionType] = []interf.Session{s}
	}

	return nil
}

func (self *SessionManager) DelSession(sessionID string) error {

	sname := strings.Split(sessionID, "/")
	if len(sname) != 2 {
		return errors.New("del session failed | invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

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

	sname := strings.Split(sessionID, "/")
	if len(sname) != 2 {
		return nil, errors.New("get session failed | invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

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
