package conn

type SessionType int

const (
	ClientSession SessionType = iota + 1
	ServerSession
)
