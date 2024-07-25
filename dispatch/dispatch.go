package dispatch

import (
	"p2p/files"
	"p2p/messages"
	"p2p/messenger"
	"p2p/transfer"
	"p2p/users"
)

type Dispatch struct {
	msger   *messenger.Messenger
	trnsfer *transfer.Transfer
}

func New(msger *messenger.Messenger) *Dispatch {
	return &Dispatch{
		msger: msger,
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
	d.msger.Send(messages.AnswerJoin(), messages.Addr(msg))
}

func (d *Dispatch) OnAnswerJoin(msg messages.Message) {
	d.msger.Send(messages.ConfirmJoin(), messages.Addr(msg))
}

func (d *Dispatch) OnConfirmJoin(msg messages.Message) {
	users.New(
		messages.Addr(msg),
		messages.UserID(msg),
		messages.Username(msg),
	)
}

func (d *Dispatch) OnUnexpected(msg messages.Message) {
	d.msger.Send(messages.BrokenProtocol(), messages.Addr(msg))
}

func (d *Dispatch) OnLocateFile(msg messages.Message) {
	key := messages.FileKey(msg)
	addr := messages.Addr(msg)
	_, found := files.Search(key)
	if !found {
		d.msger.Send(messages.FileNotFound(key), addr)
		return
	}

	d.msger.Send(messages.FileLocated(key), addr)
}

func (d *Dispatch) OnFileLocated(msg messages.Message) {
	file, found := files.Search(messages.FileKey(msg))
	if !found || file.Status != files.SEARCHING {
		return
	}
	d.trnsfer.Download(file.Key, msg)
}
