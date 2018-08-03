package conn

import (
	"errors"

	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/jumper2017/melody/errdefine"
	"github.com/jumper2017/melody/network/conn/proto"
	"google.golang.org/grpc/metadata"
)

type GrpcBidSSession struct {
	conn     proto.GrpcBid_CommServer
	recvChan chan []byte //接收数据的放置通道
	BaseSession
}

//实现 net.Addr 接口
type GBAddr struct {
	network string
	addr    string
}

func (self *GBAddr) Network() string {
	return self.network
}

func (self *GBAddr) String() string {
	return self.addr
}

func NewGBSSession(sessionType SessionType, conn proto.GrpcBid_CommServer, recvChan chan []byte) (*GrpcBidSSession, error) {

	md, ok := metadata.FromIncomingContext(conn.Context())
	if !ok {
		log.Errorf("cann't read metadata from context.")
		return nil, errdefine.ERR_CTX_PARSE
	}

	localGbaddr := new(GBAddr)
	if laddr, ok := md["server_addr"]; ok {
		localGbaddr.network = "grpc-bid"
		localGbaddr.addr = laddr[0]
	}

	remoteGbaddr := new(GBAddr)
	if raddr, ok := md["client_addr"]; ok {
		remoteGbaddr.network = "grpc-bid"
		remoteGbaddr.addr = raddr[0]
	}

	grpcBidSSession := &GrpcBidSSession{

		BaseSession: BaseSession{
			sessionType: sessionType,
			localAddr:   localGbaddr,
			remoteAddr:  remoteGbaddr,
			closeTag:    0,
		},
		conn:     conn,
		recvChan: recvChan,
	}

	return grpcBidSSession, nil
}

func (self *GrpcBidSSession) SetCbCleanSession(CbCleanSession func([]string, bool) error) error {
	if CbCleanSession != nil {
		log.Errorf("set cb clean session failed, invalid param.")
		return errors.New("set cb clean session failed, invalid param.")
	}

	self.CbCleanSession = CbCleanSession
	return nil

}

func (self *GrpcBidSSession) Start() error {

	for {

		data, err := self.conn.Recv()

		//if atomic.LoadInt32(&self.closeTag) == 1 {
		//	log.Debugf("close tag is 1, return func.")
		//	return nil
		//}

		switch {

		case err == io.EOF:
			log.Debugf("recv eof from remote.")
			self.CbCleanSession(self.needToCleanSessionsId, false)
			return nil

		case err != nil:
			log.Errorf("read data failed, err: %v", err)
			self.CbCleanSession(self.needToCleanSessionsId, true)
			continue

		}

		self.recvChan <- data.Req
	}

	return nil
}

func (self *GrpcBidSSession) Send(data []byte) error {

	rsp := &proto.GrpcBidRsp{
		Rsp: data,
	}
	return self.conn.Send(rsp)
}

//todo: 此为被动关闭
func (self *GrpcBidSSession) ClosePassive() error {

	//if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
	//	log.Debugf("close passive")
	//}
	return nil

}

func (self *GrpcBidSSession) CloseInitiative() error {
	//todo: 发送 kick 消息给客户端， 让其关闭链接

	return nil
}

func (self *GrpcBidSSession) GetSessionId() string {
	return self.sessionId
}

func (self *GrpcBidSSession) SetSessionId(sessionId string) {
	self.sessionId = sessionId
	self.needToCleanSessionsId = []string{sessionId}
	return
}
