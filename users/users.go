package users

import (
	"crypto/sha1"
	"math/big"
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

	id := new(big.Int).SetBytes(key)
	return &User{
		Addr:  addr,
		RawId: shared.HashKey(key),
		Id:    id,
	}
}
