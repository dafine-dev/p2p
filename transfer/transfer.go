package transfer

import (
	"p2p/files"
	"p2p/messages"
	"syscall"
)

type Transfer struct {
	downloadLimit int
	uploadLimit   int
	downloads     map[files.Hash]*stream
	uploads       map[files.Hash]*stream
}

func NewTransfer(downloadLimit, uploadLimit int, cancelFlags []int) *Transfer {
	return &Transfer{
		downloadLimit: downloadLimit,
		uploadLimit:   uploadLimit,
		downloads:     make(map[files.Hash]*stream),
		uploads:       make(map[files.Hash]*stream),
	}
}

func (t *Transfer) Start(port int) {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}

	addr := syscall.SockaddrInet4{
		Addr: [4]byte{0, 0, 0, 0},
		Port: port,
	}

	if err := syscall.Bind(sock, &addr); err != nil {
		panic(err)
	}

	if err := syscall.Listen(sock, syscall.SOMAXCONN); err != nil {
		panic(err)
	}

	for {
		conn, addr, err := syscall.Accept(sock)
		if err != nil {
			continue
		}

		if len(t.downloads) == t.downloadLimit {
			syscall.Close(conn)
			continue
		}

		go t.Upload(conn, addr)
	}
}

func (t *Transfer) Download(name string, msg messages.Message) *stream {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	if err = syscall.Connect(sock, msg.Addr()); err != nil {
		panic(err)
	}

	answer := make([]byte, 1024)
	n, err := syscall.Read(sock, answer)
	if err != nil {
		panic(err)
	}

	rcv_msg, ok := messages.ReadFileRequest(answer[:n])
	if !ok {
		syscall.Write(sock, messages.NewFileNotFound(rcv_msg.Key()).Raw())
		syscall.Close(sock)
		return nil
	}

	file := files.New(name)

	s := &stream{
		bufferSize: DEFAULT_BUFFER_SIZE,
		file:       file,
		stopFlag:   false,
		sock:       sock,
	}

	go s.download()
	return s
}

func (t *Transfer) Upload(sock Socket, addr syscall.Sockaddr) {

	buffer := make([]byte, 1024)

	n, err := syscall.Read(sock, buffer)
	if err != nil {
		return
	}

	msg, ok := messages.ReadFileRequest(buffer[:n])
	if !ok {
		syscall.Close(sock)
		return
	}

	key := msg.Key()
	file, found := files.Search(key)

	if !found {
		syscall.Write(sock, messages.NewFileNotFound(key).Raw())
		syscall.Close(sock)
	} else {
		s := &stream{
			file:       file,
			bufferSize: DEFAULT_BUFFER_SIZE,
			stopFlag:   false,
			sock:       sock,
		}

		t.uploads[key] = s
		go s.upload()
	}

}