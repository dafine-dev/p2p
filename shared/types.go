package shared

import (
	"math/big"
	"syscall"
)

type Socket = syscall.Handle
type Addr = syscall.SockaddrInet4
type HashKey = [20]byte
type HashId = *big.Int

var MaxId = new(big.Int).Lsh(big.NewInt(1), 160)

func Distance(a HashId, b HashId) HashId {
	switch a.Cmp(b) {
	case -1:
		// B - A
		return new(big.Int).Sub(b, a)
	case 1:
		// B + 2^N - A
		c := new(big.Int).Add(b, MaxId)
		return c.Sub(c, a)
	default:
		return new(big.Int)
	}
}

func IsBetween(a HashId, b HashId, c HashId) bool {
	newDistance := Distance(a, c)
	currentDistance := Distance(b, c)
	return newDistance.Cmp(currentDistance) < 0
}

func ReadAddr(data []byte) Addr {
	return Addr{
		Addr: [4]byte(data),
		Port: PORT,
	}
}

func ParseAddr(data ...int) Addr {
	return Addr{
		Addr: [4]byte{byte(data[0]), byte(data[1]), byte(data[2]), byte(data[3])},
		Port: PORT,
	}
}
