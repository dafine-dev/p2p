package shared

import "syscall"

const PORT = 9000

var BROADCAST_ADDR syscall.SockaddrInet4
var LOCAL_ADDR syscall.SockaddrInet4
