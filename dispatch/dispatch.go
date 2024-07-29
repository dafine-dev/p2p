package dispatch

import (
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
			break

		case messages.ANSWER_JOIN:
			go d.OnAnswerJoin(msg)
			break

		case messages.CONFIRM_JOIN:
			go d.OnConfirmJoin(msg)
			break

		case messages.INSERT_FILE:
			go d.OnInsertFile(msg)
			break

		case messages.LOCATE_FILE:
			go d.OnLocateFile(msg)
			break

		case messages.FILE_LOCATED:
			go d.OnFileLocated(msg)
			break

		default:
			go d.OnUnexpected(msg)
			break
		}
	}
}

func (d *Dispatch) OnBeginJoin(msg messages.Message) {
	d.Log(msg)
	data := messages.BeginJoin(msg)
	user := data.User()
	if d.userTable.IsSuccessor(user) {
		answer := messages.NewAnswerJoin(
			d.userTable.Current,
			d.userTable.Successor,
			d.fileTable.Between(d.userTable.Current.Id, user.Id)...,
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

	if !d.userTable.SetSuccessor(user) {
		d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), user.Addr)
		return
	}

	d.fileTable.RemoveBetween(d.userTable.Current.Id, user.Id)
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
		return
	}

	d.trnsfer.Download(data.Key(), msg)
}

func (d *Dispatch) Log(msg messages.Message) {
	log.Println(
		d.userTable.Current.Addr.Addr,
		"receives:",
		msg.Method().String(),
		"from: ",
		msg.OriginAddr().Addr)
}
