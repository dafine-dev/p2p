package messages

import (
	"syscall"
)

type Method = uint8

var LOCAL_ADDR = syscall.SockaddrInet4{
	Addr: [4]byte{0, 0, 0, 0},
	Port: 9876,
}

const (
	BEGIN_JOIN Method = iota
	ANSWER_JOIN
	CONFIRM_JOIN
	LOCATE_FILE
	REQUEST_FILE
	FILE_NOT_FOUND
	INSERT_FILE
	LEAVE
)

type Message interface {
	Addr() syscall.SockaddrInet4
	Method() Method
	Raw() []byte
}

type defaultMessage []byte

func (m defaultMessage) Addr() syscall.SockaddrInet4 {
	return syscall.SockaddrInet4{
		Addr: [4]byte(m[0:4]),
		Port: LOCAL_ADDR.Port,
	}
}

func (m defaultMessage) Method() Method {
	return m[4]
}

func (m defaultMessage) Raw() []byte {
	return m
}

func New(data []byte) defaultMessage {
	return defaultMessage(data)
}
