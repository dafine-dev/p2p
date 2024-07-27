package users

import "p2p/shared"

type UserID = [160]byte

type user struct {
	Addr shared.Addr
	Id   UserID
	Name string
}
