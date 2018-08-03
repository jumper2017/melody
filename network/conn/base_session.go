package conn

import (
	"net"

	"errors"

	log "github.com/Sirupsen/logrus"
)

type BaseSession struct {
	sessionId   string      //会话id
	sessionType SessionType //服务端还是客户端
	localAddr   net.Addr    //本端地址
	remoteAddr  net.Addr    //对端地址

	CbCleanSession func([]string, bool) error //关闭链接的回调函数，　由上层sessionManager进行操作

	closeTag int32 //关闭动作标记

	//todo: 根据业务需求进行设置，　默认在创建时设置为自身sessionId
	needToCleanSessionsId []string

	//发送和接收各一个
	//polinks package_op.PackageOpLink
	//polinkr package_op.PackageOpLink
}

func (self *BaseSession) GetSessionId() string {
	return self.sessionId
}

func (self *BaseSession) SetSessionId(sessionId string) {
	self.sessionId = sessionId
	self.needToCleanSessionsId = []string{sessionId}
	return
}

func (self *BaseSession) SetCbCleanSession(CbCleanSession func([]string, bool) error) error {
	if CbCleanSession != nil {
		log.Errorf("set cb clean session failed, invalid param.")
		return errors.New("set cb clean session failed, invalid param.")
	}

	self.CbCleanSession = CbCleanSession
	return nil

}

//func (self *BaseSession) AddPackageOp(opType package_op.PackageOpType, direct bool, params []interface{}) bool {
//	if direct {
//		return self.polinks.AddOp(opType, direct, params)
//	} else {
//		return self.polinkr.AddOp(opType, direct, params)
//	}
//}
