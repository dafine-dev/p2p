package sockets

import (
	"errors"
	"sync"
	"syscall"
)

type Socket = int

type SocketManager struct {
	Port       int
	tcpLimit   int
	udpSocket  Socket
	tcpSockets []Socket
	mutex      sync.Mutex
}

func NewManager(port int, tcpLimit int) *SocketManager {
	return &SocketManager{
		Port:       port,
		udpSocket:  -1,
		tcpLimit:   tcpLimit,
		tcpSockets: make([]Socket, 0),
	}
}

func (s *SocketManager) GetUDPSocket() (Socket, error) {
	if s.udpSocket != -1 {
		return s.udpSocket, nil
	}

	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return -1, err
	}

	return sock, nil
}

func (s *SocketManager) GetTCPSocket() (Socket, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.tcpSockets) == s.tcpLimit {
		return -1, errors.New("No avaiable sockets at the moment")
	} else if len(s.tcpSockets) == 0 {
		sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			return -1, errors.New("Error during connection creation")
		}
		return sock, nil
	} else {
		sock := s.tcpSockets[0]
		s.tcpSockets = s.tcpSockets[1:]
		return sock, nil
	}
}

func (s *SocketManager) ReleaseTCPSocket(conn Socket) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.tcpSockets = append(s.tcpSockets, conn)
}
