package interf

type Session interface {
	//SetRawConn(conn interface{}) error
	//SetRecvChan(chan interface{}) error
	Start() error

	SetCbCleanSession(func([]string) error) error

	Send(data []byte) error
	//Recv() ([]byte, error)
	Close() error

	GetSessionId() string
}
