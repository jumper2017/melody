package interf

type Session interface {
	//SetRawConn(conn interface{}) error
	//SetRecvChan(chan interface{}) error
	Start() error

	Send(data []byte) error
	//Recv() ([]byte, error)
	Close() error
}
