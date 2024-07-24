package users

import "syscall"

type UserID = [160]byte

type user struct {
	Addr syscall.SockaddrInet4
	Id   UserID
	Name string
}
