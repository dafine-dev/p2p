package main

import (
	"fmt"
	"p2p/actions"
	"p2p/shared"
	"time"
)

func main() {

	wait := make(chan struct{})
	actors := make([]*actions.Actions, 0)
	for i := 2; i < 8; i++ {
		a := actions.New(
			shared.ParseAddr(127, 0, 0, i),
			fmt.Sprintf("../p2p_test/server%d", i))
		actors = append(actors, a)
		a.Run(true)
		a.Connect()
		time.Sleep(time.Second)
	}

	time.Sleep(60 * time.Second)
	for _, actor := range actors {
		actor.PrintSuccessor()
		actor.PrintUsers()
	}

	actors[0].InsertFile("arquivo2.txt")
	actors[5].GetFile("arquivo2.txt")
	<-wait
}
