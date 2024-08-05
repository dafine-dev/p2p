package shared

import "syscall"

const PORT = 9000

var BROADCAST_ADDR syscall.SockaddrInet4 = syscall.SockaddrInet4{
	Addr: [4]byte{255, 255, 255, 255},
	Port: PORT,
}
