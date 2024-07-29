package actions

import (
	"fmt"
	"p2p/dispatch"
	"p2p/files"
	"p2p/messages"
	"p2p/messenger"
	"p2p/shared"
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
}

func New(addr shared.Addr, directory string) *Actions {
	userTable := users.StartTable(addr)
	fileManager := files.NewManager(directory)
	fileManager.SetUp()

	fmt.Println(addr)
	m := messenger.New(addr)
	f := files.NewTable()
	t := transfer.New(5, 5, fileManager, addr)
	d := dispatch.New(m, t, userTable, fileManager, f)

	return &Actions{
		msger:       m,
		userTable:   userTable,
		fileTable:   f,
		fileManager: fileManager,
		dispatcher:  d,
		trnsfer:     t,
	}
}

func (a *Actions) Run() {
	go a.msger.Run()
	go a.trnsfer.Run()
	go a.dispatcher.Run()
}

func (a *Actions) Connect() {
	addr := shared.Addr{
		Addr: [4]byte{127, 0, 0, 1},
		Port: shared.PORT,
	}
	msg := messages.NewBeginJoin(a.userTable.Current)
	for i := 2; i < 5; i++ {
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
	file := a.fileManager.Get(name)

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

}