package users

import (
	"p2p/shared"
)

type User struct {
	Addr  shared.Addr
	Id    shared.HashId
	RawId shared.HashKey
}
