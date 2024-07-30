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
			d.OnBeginJoin(msg)

		case messages.ANSWER_PRE_JOIN:
			d.OnAnswerPreJoin(msg)

		case messages.ANSWER_SUC_JOIN:
			d.OnAnswerSucJoin(msg)

		case messages.CONFIRM_JOIN:
			d.OnConfirmJoin(msg)

		case messages.INSERT_FILE:
			d.OnInsertFile(msg)

		case messages.LOCATE_FILE:
			d.OnLocateFile(msg)

		case messages.FILE_LOCATED:
			d.OnFileLocated(msg)

		case messages.FILE_NOT_FOUND:
			d.OnFileNotFound(msg)

		case messages.BROKEN_PROTOCOL:
			log.Println("Broken Protocol")

		case messages.LEAVE:
			d.OnLeave(msg)

		default:
			d.OnUnexpected(msg)

		}
	}
}

func (d *Dispatch) OnBeginJoin(msg messages.Message) {
	d.Log(msg)

	user := msg.User()
	if d.userTable.IsSuccessor(user) {
		answer := messages.NewAnswerPreJoin(d.userTable.Current)
		d.msger.Send(answer, msg.OriginAddr())
	}

	if d.userTable.IsPredecessor(user) {
		response := messages.NewAnswerSucJoin(d.userTable.Current)
		d.msger.Send(response, msg.OriginAddr())
	}
}

func (d *Dispatch) OnAnswerPreJoin(msg messages.Message) {
	d.Log(msg)

	if d.userTable.SetPredecessor(msg.User()) {
		d.msger.Send(messages.NewConfirmJoin(d.userTable.Current), msg.OriginAddr())
		return
	}

	d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), msg.OriginAddr())
}

func (d *Dispatch) OnAnswerSucJoin(msg messages.Message) {
	d.Log(msg)

	if d.userTable.SetSuccessor(msg.User()) {
		d.msger.Send(messages.NewConfirmJoin(d.userTable.Current), msg.OriginAddr())
		return
	}

	d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), msg.OriginAddr())
}

func (d *Dispatch) OnConfirmJoin(msg messages.Message) {
	d.Log(msg)
	user := msg.User()

	if d.userTable.SetSuccessor(user) {
		d.fileTable.RemoveBetween(user.Id, d.userTable.Successor.Id)
		return
	}

	d.userTable.SetPredecessor(user)

	// d.msger.Send(messages.NewBrokenProtocol(d.userTable.Current.Addr), user.Addr)
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
	log.Println(msg)
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

	user := users.New(msg.OriginAddr())
	d.userTable.Update(data.Id(), user)

	file, found := d.fileManager.Find(data.Key())
	if !found || file.Status != files.SEARCHING {
		return
	}

	d.trnsfer.Download(data.Location())
}

func (d *Dispatch) OnFileNotFound(msg messages.Message) {
	d.Log(msg)
	data := messages.FileLocated(msg)

	user := users.New(msg.OriginAddr())
	d.userTable.Update(data.Id(), user)
	file, found := d.fileManager.Find(data.Key())

	if !found || file.Status != files.SEARCHING {
		return
	}

	file.Status = files.NOT_FOUND
}

func (d *Dispatch) OnLeave(msg messages.Message) {
	d.Log(msg)

	data := messages.Leave(msg)

	if d.userTable.Successor.Id == data.User().Id {
		d.userTable.Successor = data.Neighbour()
	}

	if d.userTable.Predecessor.Id == data.User().Id {
		d.userTable.Predecessor = data.Neighbour()
	}
}

func (d *Dispatch) Log(msg messages.Message) {
	log.Println(
		d.userTable.Current.Addr.Addr,
		"receives:",
		msg.Method().String(),
		"from: ",
		msg.OriginAddr().Addr)
}
