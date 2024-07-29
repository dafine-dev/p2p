package messages

import "p2p/shared"

func NewBrokenProtocol(addr shared.Addr) Message {
	msg := make([]byte, 0)
	msg = append(msg, addr.Addr[:]...)
	msg = append(msg, byte(BROKEN_PROTOCOL))
	return msg
}
