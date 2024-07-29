package main

import (
	"p2p/actions"
	"p2p/shared"
	"time"
)

func main() {

	wait := make(chan struct{})

	s1 := actions.New(shared.ParseAddr(127, 0, 0, 2), "../p2p_test/server2")
	s1.Run()
	s1.Connect()
	time.Sleep(2 * time.Second)

	s2 := actions.New(shared.ParseAddr(127, 0, 0, 3), "../p2p_test/server3")
	s2.Run()
	s2.Connect()
	time.Sleep(2 * time.Second)

	s3 := actions.New(shared.ParseAddr(127, 0, 0, 4), "../p2p_test/server4")
	s3.Run()
	s3.Connect()

	<-wait
}
