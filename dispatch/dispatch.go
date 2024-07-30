package dispatch

import (
	"fmt"
	"log"
	"p2p/files"
	"p2p/messages"
	"p2p/messenger"
	"p2p/transfer"
	"p2p/users"
)

type Dispatch struct {
	msger       *messenger.Messenger
	trnsfer     *transfer.Transfer
	userTable   *users.Table
	fileManager *files.Manager
	fileTable   *files.Table
}

func New(msger *messenger.Messenger, t *transfer.Transfer,
	ut *users.Table, fm *files.Manager, ft *files.Table) *Dispatch {

	return &Dispatch{
		msger:       msger,
		trnsfer:     t,
		userTable:   ut,
		fileManager: fm,
		fileTable:   ft,
	}
}

func (d *Dispatch) Run() {
	for {
		msg := d.msger.Read()
		switch msg.Method() {
		case messages.BEGIN_JOIN:
			go d.OnBeginJoin(msg)

		case messages.ANSWER_JOIN:
			go d.OnAnswerJoin(msg)

		case messages.CONFIRM_JOIN:
			go d.OnConfirmJoin(msg)

		case messages.INSERT_FILE:
			go d.OnInsertFile(msg)

		case messages.LOCATE_FILE:
			go d.OnLocateFile(msg)

		case messages.FILE_LOCATED:
			go d.OnFileLocated(msg)

		case messages.FILE_NOT_FOUND:
			fmt.Println("file not found")
			go d.OnFileNotFound(msg)

		case messages.BROKEN_PROTOCOL:
			log.Println("Broken Protocol")
		default:
			go d.OnUnexpected(msg)

		}
	}
}

func (d *Dispatch) OnBeginJoin(msg messages.Message) {
	d.Log(msg)
	data := messages.BeginJoin(msg)
	user := data.User()
	if d.userTable.IsSuccessor(user) {
		b := d.fileTable.Between(user.Id, d.userTable.Successor.Id)
		// fmt.Println(b)
		answer := messages.NewAnswerJoin(
			d.userTable.Current,
			d.userTable.Successor,
			b...,
		)
		d.msger.Send(answer, data.OriginAddr())
	}
}

func (d *Dispatch) OnAnswerJoin(msg messages.Message) {
	d.Log(msg)
	data := messages.AnswerJoin(msg)
	locs, ok := data.Locations()
	if !ok {
		d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), data.OriginAddr())
		return
	}

	d.userTable.SetSuccessor(data.Successor())
	d.fileTable.Add(locs...)

	d.msger.Send(messages.NewConfirmJoin(d.userTable.Current), data.OriginAddr())
}

func (d *Dispatch) OnConfirmJoin(msg messages.Message) {
	d.Log(msg)
	data := messages.ConfirmJoin(msg)
	user := data.User()

	succ := d.userTable.Successor
	if !d.userTable.SetSuccessor(user) {
		d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), user.Addr)
		return
	}

	d.fileTable.RemoveBetween(user.Id, succ.Id)
}

func (d *Dispatch) OnInsertFile(msg messages.Message) {
	d.Log(msg)
	data := messages.InsertFile(msg)
	loc := data.Location()

	if d.userTable.Owns(loc.Id) {
		d.fileTable.Add(loc)
		return
	}

	nearest := d.userTable.Nearest(loc.Id)
	d.msger.Send(msg, nearest.Addr)
}

func (d *Dispatch) OnUnexpected(msg messages.Message) {
	d.Log(msg)
	d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), msg.OriginAddr())
}

func (d *Dispatch) OnLocateFile(msg messages.Message) {
	d.Log(msg)

	data := messages.LocateFile(msg)
	key := data.Key()
	fileId := data.Id()

	_, found := d.fileManager.Find(key)
	if found {
		response := messages.NewFileLocated(
			d.userTable.Current.Addr,
			files.NewLocation(key, d.userTable.Current.Addr))

		d.msger.Send(response, msg.OriginAddr())
		return
	}

	if d.userTable.Owns(fileId) {

		loc, found := d.fileTable.Find(key)
		if found {
			response := messages.NewFileLocated(
				d.userTable.Current.Addr,
				files.NewLocation(key, loc.Addr))
			d.msger.Send(response, msg.OriginAddr())
			return
		}

		d.msger.Send(
			messages.NewFileNotFound(d.userTable.Current.Addr, key),
			msg.OriginAddr())

		return
	}

	nearest := d.userTable.Nearest(fileId)
	d.msger.Send(msg, nearest.Addr)
}

func (d *Dispatch) OnFileLocated(msg messages.Message) {
	d.Log(msg)

	data := messages.FileLocated(msg)
	file, found := d.fileManager.Find(data.Key())
	if !found || file.Status != files.SEARCHING {
		d.userTable.Add(users.New(msg.OriginAddr()))
		return
	}

	d.trnsfer.Download(data.Location())
}

func (d *Dispatch) OnFileNotFound(msg messages.Message) {
	d.Log(msg)

	data := messages.FileLocated(msg)
	file, found := d.fileManager.Find(data.Key())

	if !found || file.Status != files.SEARCHING {
		d.userTable.Add(users.New(msg.OriginAddr()))
		return
	}

	file.Status = files.NOT_FOUND
}

func (d *Dispatch) Log(msg messages.Message) {
	log.Println(
		d.userTable.Current.Addr.Addr,
		"receives:",
		msg.Method().String(),
		"from: ",
		msg.OriginAddr().Addr)
}

func (d *Dispatch) OnLeave(msg messages.Message) {

}
