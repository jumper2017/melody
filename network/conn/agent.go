package conn

//考虑在一个　agent 中使用一个类似　leaf中　skeleton 的方式来集中处理该agent的所有事件.
//假设每个 session 中有一个chan, 用于接收到数据之后，　将数据放入该chan中，　然后agent需要
//集中监控所有session 的chan.

// session - package - processor 将一个消息的处理流程串接成一条链路， 对链路进行定制
// ! 不同的通信链路上 对于 package 的需求是不同的。
// 如何在上层灵活稳定的使用各路链接？
type Agent struct {
	sessionManager SessionManager
	recvChan       chan interface{} //从所有session接收到的数据都通过该通道传递到上层 , []byte / 解析之后的类型
}

func (self *Agent) core() {
	for {
		select {}
	}
}
