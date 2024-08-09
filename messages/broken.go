package messages

import (
	"net"
)

func NewBrokenProtocol(ip net.IP) Message {
	msg := make([]byte, 0)
	msg = append(msg, byte(BROKEN_PROTOCOL))
	msg = append(msg, ip...)
	return msg
}
