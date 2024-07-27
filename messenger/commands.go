package messenger

import (
	"p2p/messages"
	"p2p/shared"
)

type command struct {
	message  messages.Message
	destAddr shared.Socket
}
