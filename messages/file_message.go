package messages

import (
	"p2p/files"
	"p2p/shared"
)

type fileMessage struct {
	Message
}

func (f *fileMessage) Key() shared.HashKey {
	return shared.HashKey(f.Message[5])
}

func (f *fileMessage) Id() shared.HashId {
	return uint(f.Key())
}

func file_message(msg Message) *fileMessage {
	return &fileMessage{Message: msg}
}

var RequestFile = file_message
var LocateFile = file_message
var FileNotFound = file_message
var GetFile = file_message

func new_file_message(addr shared.Addr, key shared.HashKey, method Code) Message {
	msg := make([]byte, 0)
	msg = append(msg, byte(method))
	msg = append(msg, addr.Addr[:]...)
	msg = append(msg, key)
	return msg
}

func NewRequestFile(addr shared.Addr, key shared.HashKey) Message {
	return new_file_message(addr, key, REQUEST_FILE)
}

func NewLocateFile(addr shared.Addr, key shared.HashKey) Message {
	return new_file_message(addr, key, LOCATE_FILE)
}

func NewGetFile(addr shared.Addr, key shared.HashKey) Message {
	return new_file_message(addr, key, FILE)
}

func NewFileNotFound(addr shared.Addr, key shared.HashKey) Message {
	return new_file_message(addr, key, FILE_NOT_FOUND)
}

type locFile struct {
	fileMessage
}

func (f *locFile) LocationAddr() shared.Addr {
	return shared.ReadAddr(f.Message[6:10])
}

func (f *locFile) Location() *files.Location {
	return files.NewLocation(f.Key(), f.LocationAddr())
}

func loc_file(msg Message) *locFile {
	return &locFile{fileMessage: fileMessage{Message: msg}}
}

func new_loc_file(addr shared.Addr, key shared.HashKey, locationAddr shared.Addr, method Code) Message {
	msg := new_file_message(addr, key, method)
	msg = append(msg, locationAddr.Addr[:]...)
	return msg
}

func InsertFile(msg Message) *locFile {
	return loc_file(msg)
}

func FileLocated(msg Message) *locFile {
	return loc_file(msg)
}

func NewInsertFile(addr shared.Addr, loc *files.Location) Message {
	return new_loc_file(addr, loc.Key, loc.Addr, INSERT_FILE)
}

func NewFileLocated(addr shared.Addr, loc *files.Location) Message {
	return new_loc_file(addr, loc.Key, loc.Addr, FILE_LOCATED)
}
