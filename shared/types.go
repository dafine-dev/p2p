package shared

import (
	"net"
)

type Socket = int
type Addr = net.Addr
type HashKey = byte
type HashId = uint

var MaxId uint = 256

func Distance(a HashId, b HashId) HashId {
	if a < b {
		return b - a
	} else {
		return b + 256 - a
	}
}

func IsBetween(a HashId, b HashId, c HashId) bool {
	return Distance(a, c) < Distance(b, c)
}
