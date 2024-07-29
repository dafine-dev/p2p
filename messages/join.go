package messages

import (
	"p2p/files"
	"p2p/shared"
	"p2p/users"
)

type join struct {
	Message
}

func (j *join) User() *users.User {
	return users.New(j.OriginAddr())
}

type answerJoin struct {
	join
}

func (j *answerJoin) Successor() *users.User {
	return users.New(shared.ReadAddr(j.Message[5:9]))
}

func (j *answerJoin) Locations() ([]*files.Location, bool) {
	locs := make([]*files.Location, 0)
	if len(j.Message) > 9 {

		length := len(j.Message[10:])
		if length%24 != 0 {
			return locs, false
		}

		for i := 0; i < length; i += 24 {
			loc := &files.Location{
				Key:  shared.HashKey(j.Message[i : i+20]),
				Addr: shared.ReadAddr(j.Message[i+20 : i+24]),
			}
			locs = append(locs, loc)
		}
	}

	return locs, true
}

func BeginJoin(msg Message) *join {
	return &join{Message: msg}
}

func ConfirmJoin(msg Message) *join {
	return &join{Message: msg}
}

func AnswerJoin(msg Message) *answerJoin {
	return &answerJoin{
		join: join{Message: msg},
	}
}

func new_join(user *users.User, method Code) Message {
	msg := make([]byte, 0)
	msg = append(msg, byte(method))
	msg = append(msg, user.Addr.Addr[:]...)
	return msg
}

func NewBeginJoin(user *users.User) Message {
	return new_join(user, BEGIN_JOIN)
}

func new_answer_join(user, succ *users.User, method Code, locs ...*files.Location) Message {
	msg := new_join(user, method)
	msg = append(msg, succ.Addr.Addr[:]...)
	for _, loc := range locs {
		msg = append(msg, loc.Key[:]...)
		msg = append(msg, loc.Addr.Addr[:]...)
	}

	return msg
}

func NewConfirmJoin(user *users.User) Message {
	return new_join(user, CONFIRM_JOIN)
}

func NewAnswerJoin(user, succ *users.User, locs ...*files.Location) Message {
	return new_answer_join(user, succ, ANSWER_JOIN, locs...)
}
