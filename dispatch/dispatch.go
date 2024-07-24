package dispatch

import (
	"p2p/messages"
	"p2p/messenger"
)

type Dispatch struct {
	msger *messenger.Messenger
}

func New(msger *messenger.Messenger) *Dispatch {
	return &Dispatch{
		msger: msger,
	}
}

func (d *Dispatch) Run() {
	for {
		msg := d.msger.Read()
		switch msg.Method() {
		case messages.BEGIN_JOIN:
			joinMsg := messages.BeginJoin(msg)
		}
	}
}

func (d *Dispatch) OnBeginJoin(msg messages.UserMessage)
