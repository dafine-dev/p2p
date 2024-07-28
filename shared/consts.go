package shared

import "syscall"

const PORT = 9000

var BROADCAST_ADDR syscall.SockaddrInet4 = syscall.SockaddrInet4{
	Addr: [4]byte{127, 0, 0, 3},
	Port: PORT,
}
