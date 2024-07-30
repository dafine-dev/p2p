package messages

import (
	"p2p/shared"
	"p2p/users"
)

func NewBeginJoin(user *users.User) Message {
	return header(user, BEGIN_JOIN)
}

func NewConfirmJoin(user *users.User) Message {
	return header(user, CONFIRM_JOIN)
}

func NewAnswerPreJoin(user *users.User) Message {
	return header(user, ANSWER_PRE_JOIN)
}

func NewAnswerSucJoin(user *users.User) Message {
	return header(user, ANSWER_SUC_JOIN)
}

func NewLeave(user, succ *users.User) Message {
	msg := header(user, LEAVE)
	msg = append(msg, succ.Addr.Addr[:]...)

	return msg
}

type leave struct {
	Message
}

func (l *leave) Neighbour() *users.User {
	return users.New(shared.ReadAddr(l.Message[5:9]))
}

func Leave(msg Message) *leave {
	return &leave{Message: msg}
}
