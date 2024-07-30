package tracker

import (
	"math"
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
	time.Sleep(30 * time.Second)

	for i := 0; i < 7; i++ {
		id := t.UserTable.Current.Id + uint(math.Pow(2, float64(i)))
		id = id % shared.MaxId

		if t.UserTable.Owns(id) {
			t.UserTable.Update(id, t.UserTable.Current)
			continue
		}

		msg := messages.NewLocateFile(t.UserTable.Current.Addr, byte(id))
		nearest := t.UserTable.Nearest(id)
		m.Send(msg, nearest.Addr)
		// time.Sleep(5 * time.Second)
	}
}
