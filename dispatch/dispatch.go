package dispatch

import (
	"fmt"
	"log"
	"math/big"
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
		switch messages.Method(msg) {
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

	originAddr := messages.Addr(msg)
	user := messages.User(msg)

	if d.userTable.IsSuccessor(user) {
		answer := messages.AnswerJoin(
			d.userTable.Current,
			d.userTable.Successor,
			d.fileTable.Between(d.userTable.Current.Id, user.Id)...,
		)
		d.msger.Send(answer, originAddr)
	}
}

func (d *Dispatch) OnAnswerJoin(msg messages.Message) {
	d.Log(msg)
	originAddr := messages.Addr(msg)
	user := messages.User(msg)
	locs, ok := messages.FileLocations(msg)
	if !ok {
		d.msger.Send(messages.BrokenProtocol(d.userTable.Current.Addr), originAddr)
		return
	}

	d.userTable.SetSuccessor(user)
	d.fileTable.Add(locs...)

	d.msger.Send(messages.ConfirmJoin(d.userTable.Current), originAddr)
}

func (d *Dispatch) OnConfirmJoin(msg messages.Message) {
	d.Log(msg)

	originAddr := messages.Addr(msg)
	user := messages.User(msg)
	if !d.userTable.SetSuccessor(user) {
		d.msger.Send(messages.BrokenProtocol(d.userTable.Current.Addr), originAddr)
		return
	}

	d.fileTable.RemoveBetween(d.userTable.Current.Id, user.Id)
}

func (d *Dispatch) OnInsertFile(msg messages.Message) {
	d.Log(msg)

	key := messages.FileKey(msg)
	fileId := new(big.Int).SetBytes(key[:])
	loc := messages.FileLocation(msg)

	if d.userTable.Owns(fileId) {
		d.fileTable.Add(&files.Location{Addr: loc, Key: key, Id: fileId})
		return
	}

	nearest := d.userTable.Nearest(fileId)
	fmt.Println(nearest)
	d.msger.Send(msg, nearest.Addr)
}

func (d *Dispatch) OnUnexpected(msg messages.Message) {
	d.Log(msg)
	d.msger.Send(messages.BrokenProtocol(d.userTable.Current.Addr), messages.Addr(msg))
}

func (d *Dispatch) OnLocateFile(msg messages.Message) {
	key := messages.FileKey(msg)
	originAddr := messages.Addr(msg)
	fileId := new(big.Int).SetBytes(key[:])

	_, found := d.fileManager.Find(key)
	if found {
		response := messages.FileLocated(
			d.userTable.Current.Addr,
			d.userTable.Current.Addr,
			key)

		d.msger.Send(response, originAddr)
		return
	}

	if d.userTable.Owns(fileId) {

		loc, found := d.fileTable.Find(key)
		if found {
			response := messages.FileLocated(d.userTable.Current.Addr, loc.Addr, key)
			d.msger.Send(response, originAddr)
			return
		}

		d.msger.Send(messages.FileNotFound(d.userTable.Current.Addr, key), originAddr)
		return
	}

	nearest := d.userTable.Nearest(fileId)
	d.msger.Send(msg, nearest.Addr)
}

func (d *Dispatch) OnFileLocated(msg messages.Message) {
	file, found := d.fileManager.Find(messages.FileKey(msg))
	if !found || file.Status != files.SEARCHING {
		return
	}

	d.trnsfer.Download(file.Key, msg)
}

func (d *Dispatch) Log(msg messages.Message) {
	log.Println(
		d.userTable.Current.Addr.Addr,
		"receives:",
		messages.MethodName(msg),
		"from: ",
		messages.Addr(msg).Addr)
}
