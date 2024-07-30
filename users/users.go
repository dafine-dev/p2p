package users

import (
	"crypto/sha1"
	"p2p/shared"
)

type User struct {
	Addr  shared.Addr
	Id    shared.HashId
	RawId shared.HashKey
}

func New(addr shared.Addr) *User {
	hasher := sha1.New()
	hasher.Write(addr.Addr[:])
	key := hasher.Sum(nil)

	return &User{
		Addr:  addr,
		RawId: shared.HashKey(key[0]),
		Id:    uint(key[0]),
	}
}
