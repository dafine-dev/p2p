package messenger

import (
	"p2p/messages"
	"syscall"
)

type command struct {
	message  messages.Message
	destAddr syscall.SockaddrInet4
}
