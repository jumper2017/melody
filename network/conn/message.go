package conn

type Message struct {
	MessageType    uint16
	MessageId      uint16
	MessageContent interface{}
}
