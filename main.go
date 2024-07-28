package main

import (
	"fmt"
	"log"
	"p2p/dispatch"
	"p2p/files"
	"p2p/messages"
	"p2p/messenger"
	"p2p/shared"
	"p2p/transfer"
	"p2p/users"
	"time"
)

func broadcast(m *messenger.Messenger, msg messages.Message) {
	addr := shared.Addr{
		Addr: [4]byte{127, 0, 0, 1},
		Port: shared.PORT,
	}

	for i := 2; i < 5; i++ {
		addr.Addr[3] = uint8(i)
		m.Send(msg, addr)
	}
	time.Sleep(2 * time.Second)
}

func start(dir string, addr shared.Addr, filename string) {
	log.Println("Iniciando")
	userTable := users.StartTable(addr)
	fileManager := files.NewManager(dir)
	fileManager.SetUp()

	m := messenger.New(addr)
	f := files.NewTable()
	t := transfer.New(5, 5, fileManager, addr)
	d := dispatch.New(m, t, userTable, fileManager, f)
	go m.Run()
	go t.Run()
	go d.Run()

	broadcast(m, messages.BeginJoin(userTable.Current))
	// file := fileManager.Get(filename)

	go func() {

		// fmt.Println(filename)
		// fmt.Println(file)
		// fmt.Println(userTable.Current)
		time.Sleep(15 * time.Second)
		fmt.Println(userTable.Current)
		fmt.Println(userTable.Successor)
		// m.Send(messages.InsertFile(userTable.Current.Addr, file.Key), userTable.Successor.Addr)
	}()
}

func main() {

	wait := make(chan struct{})
	start("./server2", shared.Addr{Addr: [4]byte{127, 0, 0, 2}, Port: shared.PORT}, "arquivo1.txt")
	start("./server3", shared.Addr{Addr: [4]byte{127, 0, 0, 3}, Port: shared.PORT}, "arquivo4.txt")
	start("./server4", shared.Addr{Addr: [4]byte{127, 0, 0, 4}, Port: shared.PORT}, "arquivo6.txt")
	<-wait
}
