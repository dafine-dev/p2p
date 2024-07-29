package transfer

import (
	"log"
	"p2p/files"
	"p2p/messages"
	"p2p/shared"
	"syscall"
)

type Transfer struct {
	downloadLimit int
	uploadLimit   int
	downloads     map[shared.HashKey]*stream
	uploads       map[shared.HashKey]*stream
	fileManager   *files.Manager
	addr          shared.Addr
}

func New(downloadLimit, uploadLimit int,
	fileManager *files.Manager, addr shared.Addr) *Transfer {

	return &Transfer{
		downloadLimit: downloadLimit,
		uploadLimit:   uploadLimit,
		downloads:     make(map[shared.HashKey]*stream),
		uploads:       make(map[shared.HashKey]*stream),
		fileManager:   fileManager,
		addr:          addr,
	}
}

func (t *Transfer) Run() {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Println("Couldn't start TCP Socket")
		panic(err)
	}

	if err := syscall.Bind(sock, &t.addr); err != nil {
		log.Println("Couldn't bind TCP socket")
		panic(err)
	}

	if err := syscall.Listen(sock, syscall.SOMAXCONN); err != nil {
		log.Println("Couldn't listen to incoming TCP cnnections")
		panic(err)
	}

	for {
		conn, addr, err := syscall.Accept(sock)
		if err != nil {
			log.Println("Error accepting new TCP connection")
			continue
		}

		if len(t.downloads) == t.downloadLimit {
			syscall.Close(conn)
			continue
		}

		go t.Upload(conn, addr)
	}
}

func (t *Transfer) Download(key shared.HashKey, msg messages.Message) *stream {
	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	data := messages.FileLocated(msg)
	addr := data.LocationAddr()
	if err = syscall.Connect(sock, &addr); err != nil {
		syscall.Close(sock)
		return nil
	}

	_, err = syscall.Write(sock, messages.NewRequestFile(t.addr, key))
	if err != nil {
		return nil
	}

	answer := make([]byte, 1024)
	n, err := syscall.Read(sock, answer)
	if err != nil {
		return nil
	}

	if messages.Message(answer[:n]).Method() != messages.FILE {
		syscall.Write(sock, messages.NewBrokenProtocol(t.addr))
		syscall.Close(sock)
		return nil
	}

	file, _ := t.fileManager.Find(key)

	s := &stream{
		bufferSize: DEFAULT_BUFFER_SIZE,
		file:       file,
		stopFlag:   false,
		sock:       sock,
	}

	go s.download()
	return s
}

func (t *Transfer) Upload(sock shared.Socket, addr syscall.Sockaddr) {

	buffer := make([]byte, 1024)

	n, err := syscall.Read(sock, buffer)
	if err != nil {
		return
	}

	msg := buffer[:n]
	if messages.Message(msg).Method() != messages.REQUEST_FILE {
		syscall.Write(sock, messages.NewBrokenProtocol(t.addr))
		syscall.Close(sock)
		return
	}

	key := messages.RequestFile(msg).Key()
	file, found := t.fileManager.Find(key)

	if !found {
		syscall.Write(sock, messages.NewFileNotFound(t.addr, key))
		syscall.Close(sock)
		return
	}

	_, err = syscall.Write(sock, messages.NewGetFile(t.addr, key))
	if err != nil {
		return
	}

	file.Open()
	s := &stream{
		file:       file,
		bufferSize: DEFAULT_BUFFER_SIZE,
		stopFlag:   false,
		sock:       sock,
	}

	t.uploads[key] = s
	go s.upload()
}
