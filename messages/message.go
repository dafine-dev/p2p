package messages

import (
	"p2p/shared"
	"p2p/users"
)

type Code uint8

const (
	BEGIN_JOIN Code = iota
	ANSWER_PRE_JOIN
	ANSWER_SUC_JOIN
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

func (c Code) String() string {
	return map[uint8]string{
		0:  "BEGIN JOIN",
		1:  "ANSWER PRE JOIN",
		2:  "ANSWER SUC JOIN",
		3:  "CONFIRM JOIN",
		4:  "LOCATE FILE",
		5:  "REQUEST FILE",
		6:  "FILE NOT FOUND",
		7:  "FILE LOCATED",
		8:  "INSERT FILE",
		9:  "LEAVE",
		10: "FILE",
		11: "BROKEN PROTOCOL",
	}[uint8(c)]
}

type Message []byte

func (m Message) Method() Code {
	return Code(m[0])
}

func (m Message) OriginAddr() shared.Addr {
	return shared.Addr{
		Addr: [4]byte(m[1:5]),
		Port: shared.PORT,
	}
}

func (m Message) User() *users.User {
	return users.New(m.OriginAddr())
}

func header(user *users.User, method Code) Message {
	// [0] -> método
	// [1:5] -> ipv4 de origem
	msg := make(Message, 0)
	msg = append(msg, byte(method))
	msg = append(msg, user.Addr.Addr[:]...)

	return msg
}
