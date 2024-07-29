package messages

import (
	"p2p/files"
	"p2p/users"
)

func NewLeave(user, succ *users.User, locs ...*files.Location) Message {
	return new_answer_join(user, succ, LEAVE, locs...)
}

func Leave(msg Message) *answerJoin {
	return &answerJoin{
		join: join{Message: msg},
	}
}
