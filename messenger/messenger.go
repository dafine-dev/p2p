package messenger

import (
	"fmt"
	"log"
	"p2p/messages"
	"syscall"
	"time"
)

type Socket = int

type Messenger struct {
	socket    Socket
	incoming  chan messages.Message
	outcoming chan command
}

func New(port int) *Messenger {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		panic(err)
	}

	return &Messenger{
		socket:    sock,
		incoming:  make(chan messages.Message),
		outcoming: make(chan command),
	}
}

func (m *Messenger) Run() {
	go m.listenLoop()
	time.Sleep(time.Second)
	go m.writeLoop()
}

func (m *Messenger) listenLoop() {
	log.Println("Listening for UDP messages")
	for {
		buffer := make([]byte, 1024)
		n, _, err := syscall.Recvfrom(m.socket, buffer[:], 0)

		if err != nil {
			fmt.Println("Leitura de pacote UDP falhou.")
			continue
		}

		m.incoming <- messages.Message(buffer[:n])
	}
}

func (m *Messenger) writeLoop() {
	for {
		command := <-m.outcoming
		log.Println("Writing UDP message")

		err := syscall.Sendto(m.socket, command.message, 0, &command.destAddr)
		if err != nil {
			fmt.Println("Falha ao enviar mensagem", err)
			continue
		}
	}
}

func (m *Messenger) Send(msg messages.Message, to syscall.SockaddrInet4) {
	m.outcoming <- command{message: msg, destAddr: to}
}

func (m *Messenger) Read() messages.Message {
	return <-m.incoming
}
