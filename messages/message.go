package messages

import (
	"p2p/shared"
)

type Code = uint8

var LOCAL_ADDR = shared.Addr{
	Addr: [4]byte{0, 0, 0, 0},
	Port: 9876,
}

const (
	BEGIN_JOIN Code = iota
	ANSWER_JOIN
	CONFIRM_JOIN
	LOCATE_FILE
	REQUEST_FILE
	FILE_NOT_FOUND
	FILE_LOCATED
	INSERT_FILE
	LEAVE
	BROKEN_PROTOCOL
)

type Message = []byte

func Addr(msg Message) shared.Addr {
	return shared.Addr{
		Addr: [4]byte(msg[0:4]),
		Port: LOCAL_ADDR.Port,
	}
}

func Method(msg Message) Code {
	return uint8(msg[4])
}

func header(addr shared.Addr, method Code) Message {
	msg := make(Message, 0)
	msg = append(msg, LOCAL_ADDR.Addr[:]...)
	msg = append(msg, method)
	return msg
}

func BrokenProtocol() Message {
	return header(LOCAL_ADDR, BROKEN_PROTOCOL)
}
