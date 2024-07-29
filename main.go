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

func start(dir string, addr shared.Addr, filename string, flag bool, flag2 bool) {
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

	broadcast(m, messages.NewBeginJoin(userTable.Current))
	file := fileManager.Get(filename)

	if flag {

		go func() {
			time.Sleep(5 * time.Second)
			insert := messages.NewInsertFile(
				userTable.Current.Addr,
				files.NewLocation(file.Key, userTable.Current.Addr),
			)
			m.Send(insert, userTable.Successor.Addr)
		}()
	}

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println(f.All())
	}()

	if flag2 {

		go func() {
			time.Sleep(15 * time.Second)
			file := fileManager.New("arquivo4.txt")
			file.Status = files.SEARCHING
			m.Send(messages.NewLocateFile(userTable.Current.Addr, file.Key), userTable.Successor.Addr)
		}()
	}
}

func main() {

	wait := make(chan struct{})

	start("../p2p_test/server2",
		shared.Addr{Addr: [4]byte{127, 0, 0, 2}, Port: shared.PORT},
		"IWGSYUSJN.txt", false, true)

	start("../p2p_test/server3",
		shared.Addr{Addr: [4]byte{127, 0, 0, 3}, Port: shared.PORT},
		"arquivo4.txt", true, false)

	start("../p2p_test/server4",
		shared.Addr{Addr: [4]byte{127, 0, 0, 4}, Port: shared.PORT},
		"arquivo7.txt", false, false)

	<-wait
}
