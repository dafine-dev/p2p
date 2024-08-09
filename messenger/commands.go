package messenger

import (
	"net"
	"p2p/messages"
)

type command struct {
	message  messages.Message
	destAddr *net.UDPAddr
}
