package shared

import (
	"fmt"
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

func CalculateAddr() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error retrieving network interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Retrieve addresses associated with the interface
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error retrieving addresses for interface", iface.Name, ":", err)
			continue
		}

		for _, addr := range addrs {
			// Convert the address to an IP
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP

			// Skip non-IPv4 addresses
			if ip.To4() == nil {
				continue
			}

			LOCAL_ADDR = Addr{
				Addr: [4]byte(ip.To4()),
				Port: PORT,
			}

			// mask := ipnet.Mask
			// broadcastIP := make(net.IP, net.IPv4len)
			// ipv4 := ip.To4()
			// for i := range ipv4 {
			// 	broadcastIP[i] = ip[i] | ^mask[i]
			// }

			// BROADCAST_ADDR = Addr{
			// 	Addr: [4]byte(broadcastIP),
			// 	Port: PORT,
			// }
			// return
		}
	}
}
