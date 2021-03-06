package conn

import (
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/jumper2017/melody/errdefine"
	"github.com/jumper2017/melody/network/conn/proto"
	"google.golang.org/grpc/metadata"
)

type GrpcBidCSession struct {
	conn     proto.GrpcBid_CommClient
	recvChan chan []byte //接收数据的放置通道
	BaseSession
}

func NewGBCSession(sessionType SessionType, conn proto.GrpcBid_CommClient, recvChan chan []byte) (*GrpcBidCSession, error) {

	md, ok := metadata.FromIncomingContext(conn.Context())
	if !ok {
		log.Errorf("cann't read metadata from context.")
		return nil, errdefine.ERR_CTX_PARSE
	}

	localGbaddr := new(GBAddr)
	if laddr, ok := md["client_addr"]; ok {
		localGbaddr.network = "grpc-bid"
		localGbaddr.addr = laddr[0]
	}

	remoteGbaddr := new(GBAddr)
	if raddr, ok := md["server_addr"]; ok {
		remoteGbaddr.network = "grpc-bid"
		remoteGbaddr.addr = raddr[0]
	}

	grpcBidCSession := &GrpcBidCSession{
		BaseSession: BaseSession{
			sessionType: sessionType,
			localAddr:   localGbaddr,
			remoteAddr:  remoteGbaddr,
			closeTag:    0,
		},
		conn:     conn,
		recvChan: recvChan,
	}

	return grpcBidCSession, nil
}

func (self *GrpcBidCSession) Start() error {

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
			//当server端挂掉之后， 此时client recv会立即返回 err: transport is closing, 这里收到错误直接closeSend之后 退出
			log.Errorf("read data failed, err: %v", err)
			self.CbCleanSession(self.needToCleanSessionsId, true)
			return nil

		}

		self.recvChan <- data.Rsp
	}

	return nil
}

func (self *GrpcBidCSession) Send(data []byte) error {

	req := &proto.GrpcBidReq{
		Req: data,
	}
	return self.conn.Send(req)
}

//todo: 此为被动关闭
func (self *GrpcBidCSession) ClosePassive() error {

	//if atomic.CompareAndSwapInt32(&self.closeTag, 0, 1) {
	//	log.Debugf("close passive")
	//}
	return nil

}

func (self *GrpcBidCSession) CloseInitiative() error {

	return self.conn.CloseSend()
}
