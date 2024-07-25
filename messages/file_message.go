package messages

import (
	"p2p/files"
)

func FileKey(msg Message) files.Hash {
	return files.Hash(msg[5:165])
}

func fileMessage(key files.Hash, method Code) Message {
	msg := header(LOCAL_ADDR, method)
	msg = append(msg, key[:]...)
	return msg
}

func RequestFile(key files.Hash) Message {
	return fileMessage(key, REQUEST_FILE)
}

func InsertFile(name string, key files.Hash) Message {
	msg := fileMessage(key, INSERT_FILE)
	msg = append(msg, name[:]...)
	return msg
}

func FileNotFound(key files.Hash) Message {
	return fileMessage(key, FILE_NOT_FOUND)
}

func LocateFile(key files.Hash) Message {
	return fileMessage(key, LOCATE_FILE)
}

func FileLocated(key files.Hash) Message {
	return fileMessage(key, FILE_LOCATED)
}
