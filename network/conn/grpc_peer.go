package conn

import (
	"errors"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/jumper2017/melody/network/conn/proto"
	"github.com/jumper2017/melody/network/interf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////

type grpcServerImpl struct {
	sessionName    string
	recvChan       chan []byte
	funcAddSession func(sessionName string, s interf.Session)
}

func (self *grpcServerImpl) Comm(stream proto.GrpcBid_CommServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		fmt.Println("cann't read metadata from context.")
		return errors.New("read metadata failed.")
	}

	fmt.Println("get md:", md)

	grpcsSession, err := NewSession(ServerSession, stream, self.recvChan)
	if err != nil {
		logrus.Errorf("new ws session failed. err: %v", err)
		return err
	}

	self.funcAddSession(self.sessionName, grpcsSession)

	return nil
}

type GrpcPeerAcceptor struct {
	funcAddSession func(sessionName string, s interf.Session)
	listener       net.Listener
}

//接收到一个链接请求之后， 会创建一个session，
// 调用f以便将session 传入到对应的agent中
func (self *GrpcPeerAcceptor) RegisterGenerateSession(f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.funcAddSession = f
	return
}

func (self *GrpcPeerAcceptor) Start(sessionName string, listenAddr string, recvChan chan []byte) {

	var err error
	self.listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		panic("start grpc acceptor failed.")
		return
	}

	//注册自定义服务
	grpcServer := grpc.NewServer()
	proto.RegisterGrpcBidServer(grpcServer, &grpcServerImpl{
		sessionName:    sessionName,
		recvChan:       recvChan,
		funcAddSession: self.funcAddSession,
	})

	//启动grpc服务
	grpcServer.Serve(self.listener)
	return
}

func (self *GrpcPeerAcceptor) Stop() bool {
	self.listener.Close()
	return true
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
type GrpcPeerConnector struct {
	// 通过主动链接获得session之后， 存放到agent 的 session manager 中
	funcAddSession func(sessionName string, s interf.Session)
	cli            *grpc.ClientConn
}

func (self *GrpcPeerConnector) RegisterGenerateSession(f func(sessionName string, s interf.Session)) {
	if f == nil {
		panic("register generate session failed, invalid param.")
	}
	self.funcAddSession = f
	return
}

func (self *GrpcPeerConnector) Start(sessionName string, connAddr string, recvChan chan []byte) {

	//建立一个链接, 之后都是在该链接上创建流 来进行通信
	if self.cli == nil {
		var err error
		self.cli, err = grpc.Dial(connAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println("dial failed, err:", err)
			return
		}
	}

	conn := proto.NewGrpcBidClient(self.cli)

	grpccSession, err := NewSession(ClientSession, conn.(proto.GrpcBid_CommClient), recvChan)
	if err != nil {
		logrus.Errorf("new tcp session failed. err: %v", err)
		return
	}

	go self.funcAddSession(sessionName, grpccSession)

	return
}
