package network

//定义peer接口，　表示的是各种协议的connector/acceptor，
//对于　connector 来说，　需要提供一个创建原始链接的接口
//对于　acceptor　来说，　需要提供一个启动监听的接口

//peer的功能主要包含如下：
//1. 原始链接的建立还删除 考虑建立一个子接口，　各个协议分别实现子接口，　包括各协议相关属性
//2. 基于原始链接来创建session和管理众多session，　考虑建立一个子接口，　各个协议分别实现子接口
//3. 是否需要单独定义错误处理策略，　是否需要就不同协议分别定义

//定义session接口, 表示的是各种协议的会话抽象
//session的功能主要包含如下：
//1. 通过原始链接进行数据收发, 数据接收肯定是在一个单独协程, 数据处理和发送在同一个协程中, 即异步接收,同步发送.

//?这里需要考虑一个问题： tcp/ws/rpc-async/udp/http 协议都可以抽象出 read/write 接口, 但是如何对 rpc-sync/timer/db 进行抽象
//考虑两种方式：
//1. 定义一个接口

// 考虑其由多个部分来组成：
//1. 协议相关属性，　2. 协议相关通信过程（connect / accept），　3.　错误处理策略（发送/接受错误）

//在定义一个系统的时候, 是否能够分解成两个部分:　定义模块　定义模块的链接方式． 比如这里可以定义从通信层到逻辑层的链接方式为channel
//定义从逻辑层到通信层的链接方式为　simple (即直接调用)

import (
	"errors"
	"strconv"
	"strings"
	"sync"
)

//PeerAcceptor 获取服务端session
//PeerConnector 获取客户端session
//SessionManager 用于保存一个用户对象的所有session
//Agent 用于表示一个用户对象， 包含该用户的session manager

type PeerAcceptor interface {
	Start(addr string)
	Stop() bool
	//接收到一个链接请求之后， 会创建一个session，
	// 调用f以便将session 传入到对应的agent中
	RegisterGenerateSession(f func(s Session))
}

type PeerConnector interface {
	// 通过主动链接获得session之后， 存放到agent 的 session manager 中
	Start(addr string) Session
}

type SessionManagerOps interface {
	AddSession(sessionType string, s Session) error
	// sessionID 为 "tpc:localip:localport:remoteip:remoteport/1" 表示删除 tpc:localip:localport:remoteip:remoteport 的第二个元素，
	// "tpc:localip:localport:remoteip:remoteport/*" 表示删除 tpc:localip:localport:remoteip:remoteport 的所有元素
	DelSession(sessionID string) error
	GetSession(sessionID string) ([]Session, error) // 同 DelSession 参数含义， 返回可能是一个或多个Session
	//CloseSession(sessionID string) error //同 DelSession 参数含义
}

type SessionManager struct {
	// 使用字符串来作为 sessionType, 由三个部分组成 tpc:localip:localport:remoteip:remoteport
	// 这里考虑同一类session可以存在多个数量， "tpc:localip:localport:remoteip:remoteport/0"表示sessionID, 即
	// 链接到gate的tcp协议类型的index 0 的session, 即  对应value 的第一个元素
	sessions     map[string][]Session
	sessionsLock sync.Mutex
	SessionManagerOps
}

func (self *SessionManager) AddSession(sessionType string, s Session) error {

	if s == nil {
		return errors.New("add session failed | invalid param.")
	}

	self.sessionsLock.Lock()
	defer self.sessionsLock.Unlock()

	if slice, ok := self.sessions[sessionType]; ok {
		slice = append(slice, s)
	} else {
		self.sessions[sessionType] = []Session{s}
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

func (self *SessionManager) GetSession(sessionID string) ([]Session, error) {

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
		return []Session{self.sessions[sname[0]][index]}, nil
	}

	return nil, errors.New("get session failed | invalid param.")
}

type Agent struct {
	sessionManager SessionManager
}
