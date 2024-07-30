package messages

import (
	"fmt"
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

type leave struct {
	join
}

func (j *leave) Successor() *users.User {
	return users.New(shared.ReadAddr(j.Message[5:9]))
}

func (j *leave) Locations() ([]*files.Location, bool) {
	locs := make([]*files.Location, 0)
	if len(j.Message) > 9 {

		length := len(j.Message[9:])
		if length%5 != 0 {
			return locs, false
		}

		for i := 9; i < length; i += 5 {
			loc := &files.Location{
				Key:  shared.HashKey(j.Message[i]),
				Addr: shared.ReadAddr(j.Message[i+1 : i+5]),
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
	fmt.Println(msg)
	return &join{Message: msg}
}

func Leave(msg Message) *leave {
	return &leave{
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

func NewLeave(user, succ *users.User, locs ...*files.Location) Message {
	msg := new_join(user, LEAVE)
	msg = append(msg, succ.Addr.Addr[:]...)
	for _, loc := range locs {
		msg = append(msg, loc.Key)
		msg = append(msg, loc.Addr.Addr[:]...)
	}

	return msg
}

func NewConfirmJoin(user *users.User) Message {
	return new_join(user, CONFIRM_JOIN)
}

func NewAnswerSucJoin(user *users.User) Message {
	return new_join(user, ANSWER_SUC_JOIN)
}

type answerPreJoin struct {
	join
}

func (j *answerPreJoin) Locations() ([]*files.Location, bool) {
	locs := make([]*files.Location, 0)
	if len(j.Message) > 5 {

		length := len(j.Message[5:])
		if length%5 != 0 {
			return locs, false
		}

		for i := 5; i < length; i += 5 {
			loc := &files.Location{
				Key:  shared.HashKey(j.Message[i]),
				Addr: shared.ReadAddr(j.Message[i+1 : i+5]),
			}
			locs = append(locs, loc)
		}
	}
	return locs, true
}

func AnswerSucJoin(msg Message) *join {
	return &join{Message: msg}
}

func AnswerPreJoin(msg Message) *answerPreJoin {
	return &answerPreJoin{join: join{Message: msg}}
}

func NewAnswerPreJoin(user *users.User, locs ...*files.Location) Message {
	msg := new_join(user, ANSWER_PRE_JOIN)
	for _, loc := range locs {
		msg = append(msg, loc.Key)
		msg = append(msg, loc.Addr.Addr[:]...)
	}

	return msg
}
