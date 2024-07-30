package actions

import (
	"fmt"
	"p2p/dispatch"
	"p2p/files"
	"p2p/messages"
	"p2p/messenger"
	"p2p/shared"
	"p2p/tracker"
	"p2p/transfer"
	"p2p/users"
)

type Actions struct {
	msger       *messenger.Messenger
	userTable   *users.Table
	fileTable   *files.Table
	fileManager *files.Manager
	trnsfer     *transfer.Transfer
	dispatcher  *dispatch.Dispatch
	trcker      *tracker.Tracker
}

func New(addr shared.Addr, directory string) *Actions {
	userTable := users.StartTable(addr)
	fileManager := files.NewManager(directory)
	fileManager.SetUp()

	m := messenger.New(addr)
	f := files.NewTable()
	t := transfer.New(15, 15, fileManager, addr)
	d := dispatch.New(m, t, userTable, fileManager, f)
	tr := &tracker.Tracker{UserTable: userTable}

	return &Actions{
		msger:       m,
		userTable:   userTable,
		fileTable:   f,
		fileManager: fileManager,
		dispatcher:  d,
		trnsfer:     t,
		trcker:      tr,
	}
}

func (a *Actions) Run(tracking bool) {
	go a.msger.Run()
	go a.trnsfer.Run()
	go a.dispatcher.Run()
	if tracking {
		go a.trcker.Run(a.msger)
	}
}

func (a *Actions) Connect() {
	addr := shared.Addr{
		Addr: [4]byte{127, 0, 0, 1},
		Port: shared.PORT,
	}
	msg := messages.NewBeginJoin(a.userTable.Current)
	for i := 2; i < 18; i++ {
		addr.Addr[3] = uint8(i)
		a.msger.Send(msg, addr)
	}
	// a.msger.Send(messages.NewBeginJoin(a.userTable.Current), shared.BROADCAST_ADDR)
}

func (a *Actions) InsertFile(name string) {
	addr := a.userTable.Current.Addr
	file := a.fileManager.Get(name)
	loc := files.NewLocation(file.Key, addr)

	if a.userTable.Owns(file.Id) {
		a.fileTable.Add(loc)
		return
	}

	nearest := a.userTable.Nearest(file.Id)
	a.msger.Send(messages.NewInsertFile(addr, loc), nearest.Addr)
}

func (a *Actions) GetFile(name string) {
	file := a.fileManager.New(name)

	loc, found := a.fileTable.Find(file.Key)
	if found {
		a.trnsfer.Download(loc)
		return
	}

	nearest := a.userTable.Nearest(file.Id)
	file.Status = files.SEARCHING
	a.msger.Send(messages.NewLocateFile(a.userTable.Current.Addr, file.Key), nearest.Addr)
}

func (a *Actions) Leave() {
	user := a.userTable.Current
	succ := a.userTable.Successor
	pred := a.userTable.Predecessor

	msg := messages.NewLeave(user, succ,
		a.fileTable.Between(user.Id, succ.Id)...)

	a.msger.Send(msg, pred.Addr)

	msg = messages.NewLeave(user, pred)
	a.msger.Send(msg, succ.Addr)
}

func (a *Actions) FileTable() map[shared.HashKey]*files.Location {
	return a.fileTable.All()
}

func (a *Actions) PrintSuccessor() {
	fmt.Println(
		a.userTable.Predecessor.Id,
		a.userTable.Current.Id,
		a.userTable.Successor.Id)
}

func (a *Actions) PrintUsers() {
	m := make([]string, 0)
	for key, user := range a.userTable.All() {
		m = append(m, fmt.Sprintf("| %d:%d |", user.Id, key))
	}

	fmt.Println(m)
}
