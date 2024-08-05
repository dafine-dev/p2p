package shared

import (
	"net"
	"syscall"
)

type Socket = int
type Addr = syscall.SockaddrInet4
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

func LocalAddr() Addr {
	var a Addr
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			a = Addr{
				Addr: [4]byte(ip.To4()),
				Port: PORT,
			}
		}
	}

	return a
}
