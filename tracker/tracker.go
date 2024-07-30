package tracker

import (
	"fmt"
	"math/big"
	"p2p/messages"
	"p2p/messenger"
	"p2p/shared"
	"p2p/users"
	"time"
)

type Tracker struct {
	UserTable *users.Table
}

func (t *Tracker) Run(m *messenger.Messenger) {
	time.Sleep(7 * time.Second)
	// i := uint(0)
	// id := new(big.Int)
	// for {
	// 	p := new(big.Int).Lsh(big.NewInt(2), i)
	// 	id.Add(t.UserTable.Current.Id, p)

	// 	nearest := t.UserTable.Nearest(id)
	// 	msg := messages.NewLocateFile(t.UserTable.Current.Addr, [20]byte(id.Bytes()))
	// 	m.Send(msg, nearest.Addr)

	// 	time.Sleep(1 * time.Second)

	// 	if i == 2 {
	// 		break
	// 	} else {
	// 		i++
	// 	}

	// }

	two := big.NewInt(2)
	id := new(big.Int)
	two.Exp(two, big.NewInt(156), nil)
	id.Add(t.UserTable.Current.Id, two)
	fmt.Println(id.String())
	id.Mod(id, shared.MaxId)
	fmt.Println(id.String())
	var key shared.HashKey
	value := id.Bytes()
	copy(key[20-len(value):], id.Bytes())
	msg := messages.NewLocateFile(t.UserTable.Current.Addr, key)

	nearest := t.UserTable.Nearest(id)
	flag := t.UserTable.Owns(id)
	if !flag {
		m.Send(msg, nearest.Addr)
	}

	fmt.Println(flag)
}
