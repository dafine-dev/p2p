package main

import (
	"log"
	"p2p/dispatch"
	"p2p/messages"
	"p2p/messenger"
	"p2p/transfer"
	"syscall"
)

const PORT = 8080

func main() {
	log.Println("Iniciando")
	wait := make(chan struct{})

	m := messenger.New(PORT)
	t := transfer.New(5, 5)
	d := dispatch.New(m)
	go m.Run()
	go t.Start(PORT)
	go d.Run()

	addr := syscall.SockaddrInet4{
		Addr: [4]byte{255, 255, 255, 255},
		Port: PORT,
	}
	m.Send(messages.BeginJoin(), addr)
	<-wait
}
