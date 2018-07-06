package conn

//考虑在一个　agent 中使用一个类似　leaf中　skeleton 的方式来集中处理该agent的所有事件.
//假设每个 session 中有一个chan, 用于接收到数据之后，　将数据放入该chan中，　然后agent需要
//集中监控所有session 的chan.
type Agent struct {
	sessionManager SessionManager
	recvChan       chan interface{} //从所有session接收到的数据都通过该通道传递到上层 , []byte / 解析之后的类型
}

func (self *Agent) core() {
	for {
		select {}
	}
}
