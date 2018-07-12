package conn

//流式消息中的结构　length+content
const STREAM_MSG_HEAD_LENGTH = 2

//session类型
type SessionType int

const (
	ClientSession SessionType = iota + 1
	ServerSession
)
