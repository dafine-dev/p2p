package messages

import (
	"p2p/shared"
)

type Code = uint8

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
	FILE
	BROKEN_PROTOCOL
)

func MethodName(msg Message) string {
	return map[uint8]string{
		0:  "BEGIN JOIN",
		1:  "ANSWER JOIN",
		2:  "CONFIRM JOIN",
		3:  "LOCATE FILE",
		4:  "REQUEST FILE",
		5:  "FILE NOT FOUND",
		6:  "FILE LOCATED",
		7:  "INSERT FILE",
		8:  "LEAVE",
		9:  "FILE",
		10: "BROKEN PROTOCOL",
	}[Method(msg)]
}

type Message = []byte

func Addr(msg Message) shared.Addr {
	return shared.Addr{
		Addr: [4]byte(msg[0:4]),
		Port: shared.PORT,
	}
}

func Method(msg Message) Code {
	return uint8(msg[4])
}

func header(addr shared.Addr, method Code) Message {
	msg := make(Message, 0)
	msg = append(msg, addr.Addr[:]...)
	msg = append(msg, method)
	return msg
}

func BrokenProtocol(addr shared.Addr) Message {
	return header(addr, BROKEN_PROTOCOL)
}
