package messenger

import (
	"fmt"
	"log"
	"p2p/messages"
	"p2p/shared"
	"syscall"
)

type Messenger struct {
	socket    shared.Socket
	incoming  chan messages.Message
	outcoming chan command
	addr      shared.Addr
}

func New(addr shared.Addr) *Messenger {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Println("Couldn't start")
		panic(err)
	}

	err = syscall.SetsockoptInt(sock, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		log.Println("Couldn't enable broadcast in UDP socket.")
		panic(err)
	}

	syscall.Bind(sock, &addr)

	return &Messenger{
		socket:    sock,
		incoming:  make(chan messages.Message),
		outcoming: make(chan command),
		addr:      addr,
	}
}

func (m *Messenger) Run() {
	go m.listenLoop()
	go m.writeLoop()
}

func (m *Messenger) listenLoop() {
	for {
		buffer := make([]byte, 1024)
		n, addr, err := syscall.Recvfrom(m.socket, buffer[:], 0)

		if err != nil {
			panic(err)
		}

		addrI4, ok := addr.(*syscall.SockaddrInet4)
		if !ok {
			continue
		}

		if addrI4.Addr == m.addr.Addr {
			continue
		}

		m.incoming <- messages.Message(buffer[:n])
	}
}

func (m *Messenger) writeLoop() {
	for {
		command := <-m.outcoming
		// log.Println("Writing UDP message")

		log.Println(command.destAddr)
		err := syscall.Sendto(m.socket, command.message, 0, &command.destAddr)
		if err != nil {
			fmt.Println("Falha ao enviar mensagem", err)
			continue
		}
	}
}

func (m *Messenger) Send(msg messages.Message, to shared.Addr) {
	m.outcoming <- command{message: msg, destAddr: to}
}

func (m *Messenger) Read() messages.Message {
	return <-m.incoming
}
