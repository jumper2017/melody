package conn

type Message struct {
	MessageType    uint16
	MessageId      uint16
	MessageContent interface{}
}

//1. session 使用同一个 chan []byte 将数据上传到解码层， 解码层需要获得是哪个session，然后使用对应的 parser 进行解码
//可以考虑在 上传的数据前面加几个字节记录 sessionId ，这样 parser层就知道是哪个 session ， 而已获得对应的 parser实例进行解析
//2. parser 实例考虑参考下 leaf
//对于parser的定义 还关系到 message 的定义。 考虑直接采用 leaf 中的方式, 高度抽象 message,
