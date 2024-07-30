package messages

import "p2p/shared"

func NewBrokenProtocol(addr shared.Addr) Message {
	msg := make([]byte, 0)
	msg = append(msg, byte(BROKEN_PROTOCOL))
	msg = append(msg, addr.Addr[:]...)
	return msg
}
