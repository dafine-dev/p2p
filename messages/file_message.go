package messages

import "p2p/files"

type fileMessage struct {
	defaultMessage
}

func (f fileMessage) Key() files.Hash {
	return files.Hash(f.defaultMessage[5:165])
}

func newFileMessage(key files.Hash, method Method) fileMessage {
	data := make([]byte, 0)
	data = append(data, LOCAL_ADDR.Addr[:]...)
	data = append(data, method)
	data = append(data, key[:]...)
	return fileMessage{
		defaultMessage: data,
	}
}

func NewFileRequest(key files.Hash) fileMessage {
	return newFileMessage(key, FILE_REQUEST)
}

func ReadFileRequest(data []byte) (fileMessage, bool) {
	if data[4] == FILE_REQUEST {
		return fileMessage{}, false
	} else {
		return fileMessage{defaultMessage: data}, true
	}
}

func NewFileNotFound(key files.Hash) fileMessage {
	return newFileMessage(key, FILE_NOT_FOUND)
}

func ReadFileNotFound(data []byte) (fileMessage, bool) {
	if data[4] == FILE_NOT_FOUND {
		return fileMessage{}, false
	} else {
		return fileMessage{defaultMessage: data}, true
	}
}
