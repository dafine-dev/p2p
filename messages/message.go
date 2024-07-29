package messages

import "p2p/shared"

type Code uint8

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

func (c Code) String() string {
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
