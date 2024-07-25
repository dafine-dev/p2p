package messages

import (
	"crypto/sha1"
	"p2p/users"
)

func UserID(msg Message) users.UserID {
	return users.UserID(msg[5:165])
}

func Username(msg Message) string {
	return string(msg[165 : 165+256])
}

func join(name string, method Code) Message {
	msg := header(LOCAL_ADDR, method)

	hasher := sha1.New()
	hasher.Write(LOCAL_ADDR.Addr[:])

	msg = append(msg, hasher.Sum(nil)...)
	msg = append(msg, []byte(name)...)
	return msg
}

func BeginJoin() Message {
	return join(users.CURRENT_USERNAME, BEGIN_JOIN)
}

func AnswerJoin() Message {
	return join(users.CURRENT_USERNAME, ANSWER_JOIN)
}

func ConfirmJoin() Message {
	return header(LOCAL_ADDR, CONFIRM_JOIN)
}

func Leave() Message {
	return header(LOCAL_ADDR, LEAVE)
}
