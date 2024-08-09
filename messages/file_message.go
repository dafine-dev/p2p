package messages

import (
	"net"
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

func new_file_message(ip net.IP, key shared.HashKey, method Code) Message {
	msg := make([]byte, 0)
	msg = append(msg, byte(method))
	msg = append(msg, ip...)
	msg = append(msg, key)
	return msg
}

func NewRequestFile(ip net.IP, key shared.HashKey) Message {
	return new_file_message(ip, key, REQUEST_FILE)
}

func NewLocateFile(ip net.IP, key shared.HashKey) Message {
	return new_file_message(ip, key, LOCATE_FILE)
}

func NewGetFile(ip net.IP, key shared.HashKey) Message {
	return new_file_message(ip, key, FILE)
}

func NewFileNotFound(ip net.IP, key shared.HashKey) Message {
	return new_file_message(ip, key, FILE_NOT_FOUND)
}

type locFile struct {
	fileMessage
}

func (f *locFile) LocationAddr() net.IP {
	return net.IP(f.Message[6:10])
}

func (f *locFile) Location() *files.Location {
	return files.NewLocation(f.Key(), f.LocationAddr())
}

func loc_file(msg Message) *locFile {
	return &locFile{fileMessage: fileMessage{Message: msg}}
}

func new_loc_file(ip net.IP, key shared.HashKey, locationIP net.IP, method Code) Message {
	msg := new_file_message(ip, key, method)
	msg = append(msg, locationIP...)
	return msg
}

func InsertFile(msg Message) *locFile {
	return loc_file(msg)
}

func FileLocated(msg Message) *locFile {
	return loc_file(msg)
}

func NewInsertFile(ip net.IP, loc *files.Location) Message {
	return new_loc_file(ip, loc.Key, loc.IP, INSERT_FILE)
}

func NewFileLocated(ip net.IP, loc *files.Location) Message {
	return new_loc_file(ip, loc.Key, loc.IP, FILE_LOCATED)
}
