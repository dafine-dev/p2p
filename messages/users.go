package messages

import "p2p/users"

type UserMessage struct {
	defaultMessage
}

func (u *UserMessage) Id() users.UserID {
	return [160]byte(u.defaultMessage.Raw()[6:167])
}

func (u *UserMessage) Name() string {
	return string(u.defaultMessage.Raw()[167 : 167+256])
}
