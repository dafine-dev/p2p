package messages

import (
	"math/big"
	"p2p/files"
	"p2p/shared"
	"p2p/users"
)

func UserId(msg Message) shared.HashId {
	return new(big.Int).SetBytes(msg[5:25])
}

func RawId(msg Message) shared.HashKey {
	return shared.HashKey(msg[5:26])
}

func join(user *users.User, method Code) Message {
	msg := header(user.Addr, method)

	msg = append(msg, user.RawId[:]...)
	return msg
}

func BeginJoin(user *users.User) Message {
	msg := join(user, BEGIN_JOIN)
	msg = append(msg, user.RawId[:]...)
	return msg
}

func AnswerJoin(user *users.User, suc *users.User, locs ...*files.Location) Message {
	msg := join(user, ANSWER_JOIN)
	msg = append(msg, user.Addr.Addr[:]...)
	msg = append(msg, user.RawId[:]...)
	for _, loc := range locs {
		msg = append(msg, loc.Key[:]...)
	}
	return msg
}

func ConfirmJoin(user *users.User) Message {
	return header(user.Addr, CONFIRM_JOIN)
}

func Leave(user *users.User) Message {
	return header(user.Addr, LEAVE)
}

func User(msg Message) *users.User {
	return &users.User{
		Addr:  Addr(msg),
		RawId: RawId(msg),
		Id:    UserId(msg),
	}
}

func FileLocations(msg Message) ([]*files.Location, bool) {
	data := msg[49:]
	locs := make([]*files.Location, 0)
	if len(data)%24 != 0 {
		return locs, false
	}

	for i := 0; i < len(data); i += 24 {
		key := msg[i : i+20]

		loc := &files.Location{
			Key: shared.HashKey(key),
			Id:  new(big.Int).SetBytes(key),
			Addr: shared.Addr{
				Addr: [4]byte(msg[i+20:]),
				Port: shared.PORT,
			}}
		locs = append(locs, loc)
	}

	return locs, true
}
