package messages

import (
	"p2p/shared"
)

func FileKey(msg Message) shared.HashKey {
	return shared.HashKey(msg[5:25])
}

func FileLocation(msg Message) shared.Addr {
	return shared.Addr{
		Addr: [4]byte(msg[25:29]),
		Port: shared.PORT,
	}
}

func fileMessage(addr shared.Addr, key shared.HashKey, method Code) Message {
	msg := header(addr, method)
	msg = append(msg, key[:]...)
	return msg
}

func RequestFile(addr shared.Addr, key shared.HashKey) Message {
	return fileMessage(addr, key, REQUEST_FILE)
}

func InsertFile(addr shared.Addr, key shared.HashKey) Message {
	return fileMessage(addr, key, INSERT_FILE)
}

func FileNotFound(addr shared.Addr, key shared.HashKey) Message {
	return fileMessage(addr, key, FILE_NOT_FOUND)
}

func LocateFile(addr shared.Addr, key shared.HashKey) Message {
	return fileMessage(addr, key, LOCATE_FILE)
}

func FileLocated(addr shared.Addr, locAddr shared.Addr, key shared.HashKey) Message {
	msg := fileMessage(addr, key, FILE_LOCATED)
	msg = append(msg, locAddr.Addr[:]...)
	return msg
}

func File(addr shared.Addr, key shared.HashKey) Message {
	return fileMessage(addr, key, FILE)
}
