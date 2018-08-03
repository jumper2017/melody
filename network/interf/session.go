package interf

type Session interface {
	//SetRawConn(conn interface{}) error
	//SetRecvChan(chan interface{}) error
	Start() error

	SetCbCleanSession(func([]string, bool) error) error

	Send(data []byte) error
	//Recv() ([]byte, error)
	ClosePassive() error
	CloseInitiative() error

	GetSessionId() string
	SetSessionId(string)
}
