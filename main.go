package main

import (
	"p2p/actions"
	"p2p/shared"
	"time"
)

func main() {

	wait := make(chan struct{})
	s2 := actions.New(shared.ParseAddr(127, 0, 0, 2), "../p2p_test/server2")
	s2.Run(true)
	s2.Connect()
	time.Sleep(2 * time.Second)

	s3 := actions.New(shared.ParseAddr(127, 0, 0, 3), "../p2p_test/server3")
	s3.Run(false)
	s3.Connect()
	time.Sleep(2 * time.Second)

	s4 := actions.New(shared.ParseAddr(127, 0, 0, 4), "../p2p_test/server4")
	s4.Run(false)
	s4.Connect()
	time.Sleep(5 * time.Second)

	// s2.PrintSuccessor()
	// s3.PrintSuccessor()
	// s4.PrintSuccessor()
	// fmt.Println("s3", s3.FileTable())
	// s2.InsertFile("arquivo1.txt")
	// s4.InsertFile("arquivo9.txt")
	// time.Sleep(1 * time.Second)

	// fmt.Println("s2", s2.FileTable())
	// fmt.Println("s4", s4.FileTable())

	// time.Sleep(2 * time.Second)

	// s2.GetFile("arquivo9.txt")

	// fmt.Println("s2", s2.FileTable())
	// fmt.Println("s3", s3.FileTable())
	// fmt.Println("s4", s4.FileTable())
	// s2.PrintSuccessor()
	// s3.PrintSuccessor()
	// s4.PrintSuccessor()
	s2.PrintUsers()
	<-wait
}
