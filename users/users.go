package users

import (
	"crypto/sha1"
	"net"
	"p2p/shared"
)

type User struct {
	IP    net.IP
	Id    shared.HashId
	RawId shared.HashKey
}

func New(ip net.IP) *User {
	hasher := sha1.New()
	hasher.Write(ip)
	key := hasher.Sum(nil)

	return &User{
		IP:    ip,
		RawId: shared.HashKey(key[0]),
		Id:    uint(key[0]),
	}
}
