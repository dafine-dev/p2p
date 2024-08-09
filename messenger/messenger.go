package messenger

import (
	"log"
	"net"
	"p2p/messages"
	"p2p/shared"
)

type Messenger struct {
	conn      *net.UDPConn
	incoming  chan messages.Message
	outcoming chan command
	ip        net.IP
}

func New(ip net.IP) *Messenger {
	addr := net.UDPAddr{
		IP:   ip,
		Port: shared.PORT,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println("Couldn't start UDP listening socket.")
		return nil
	}

	return &Messenger{
		conn:      conn,
		incoming:  make(chan messages.Message),
		outcoming: make(chan command),
	}
}

func (m *Messenger) Run() {
	go m.listenLoop()
	go m.writeLoop()
}

func (m *Messenger) listenLoop() {
	for {
		buffer := make([]byte, 1024)
		n, addr, err := m.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Couldn't read message from UDP socket.")
			continue
		}

		if addr.IP.Equal(m.ip) {
			continue
		}

		m.incoming <- messages.Message(buffer[:n])
	}
}

func (m *Messenger) writeLoop() {
	for {
		command := <-m.outcoming
		// log.Println("Writing UDP message")

		_, err := m.conn.WriteTo(command.message, command.destAddr)
		if err != nil {
			log.Printf("Failed to send message %s to %s\n",
				command.message.Method().String(),
				command.destAddr.String())
			continue
		}
	}
}

func (m *Messenger) Send(msg messages.Message, to net.IP) {
	addr := net.UDPAddr{
		IP:   to,
		Port: shared.PORT,
	}
	m.outcoming <- command{message: msg, destAddr: &addr}
}

func (m *Messenger) Read() messages.Message {
	return <-m.incoming
}
