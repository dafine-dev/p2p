package shared

import (
	"math/big"
	"syscall"
)

type Socket = int
type Addr = syscall.SockaddrInet4
type HashKey = [20]byte
type HashId = *big.Int

var MaxId = new(big.Int).Exp(big.NewInt(2), big.NewInt(160), nil)

func Distance(a HashId, b HashId) HashId {
	r := a.Cmp(b)
	if r < 0 {
		return new(big.Int).Sub(b, a)
	} else {
		c := new(big.Int).Add(b, MaxId)
		return c.Sub(c, a)
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
