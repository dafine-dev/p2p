package transfer

import (
	"log"
	"net"
	"p2p/files"
	"p2p/messages"
	"p2p/shared"
)

type Transfer struct {
	downloadLimit int
	uploadLimit   int
	downloads     map[shared.HashKey]*stream
	uploads       map[shared.HashKey]*stream
	fileManager   *files.Manager
	addr          *net.TCPAddr
}

func New(downloadLimit, uploadLimit int,
	fileManager *files.Manager, ip net.IP) *Transfer {

	addr := net.TCPAddr{
		IP:   ip,
		Port: shared.PORT,
	}
	return &Transfer{
		downloadLimit: downloadLimit,
		uploadLimit:   uploadLimit,
		downloads:     make(map[shared.HashKey]*stream),
		uploads:       make(map[shared.HashKey]*stream),
		fileManager:   fileManager,
		addr:          &addr,
	}
}

func (t *Transfer) Run() {
	listener, err := net.ListenTCP("tcp", t.addr)

	if err != nil {
		log.Println("Couldn't start File transfering server. TCP socket failed to open.")
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Error establising TCP connection from %s\n", conn.RemoteAddr())
		}

		if len(t.downloads) == t.downloadLimit {
			conn.Close()
			continue
		}

		go t.Upload(conn)
	}
}

func (t *Transfer) Download(loc *files.Location) *stream {
	srcAddr := net.TCPAddr{
		IP:   loc.IP,
		Port: shared.PORT,
	}
	conn, err := net.DialTCP("tcp", nil, &srcAddr)
	if err != nil {
		log.Println("Couldn't start download of file with key %s", loc.Key)
		return nil
	}

	n, err := conn.Write(messages.NewRequestFile(t.addr.IP, loc.Key))
	if err != nil {
		log.Println("Couldn't proceed with file donwload. Target rejected file request.")
		conn.Close()
		return nil
	}

	answer := make([]byte, 1024)
	n, err = conn.Read(answer)
	if err != nil {
		log.Println("Error reading target response to file request. Canceling download.")
		conn.Close()
		return nil
	}

	if messages.Message(answer[:n]).Method() != messages.FILE {
		conn.Write(messages.NewBrokenProtocol(t.addr.IP))
		conn.Close()
		log.Println("Target answered with wrong message. Closing connection.")
		return nil
	}

	file, _ := t.fileManager.Find(loc.Key)

	s := &stream{
		bufferSize: DEFAULT_BUFFER_SIZE,
		file:       file,
		stopFlag:   false,
		conn:       conn,
	}

	go s.download()
	return s
}

func (t *Transfer) Upload(conn *net.TCPConn) {
	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Couldn't read from incoming TCP connection. Closing it.")
		conn.Close()
		return
	}

	msg := buffer[:n]
	if messages.Message(msg).Method() != messages.REQUEST_FILE {
		conn.Write(messages.NewBrokenProtocol(t.addr.IP))
		conn.Close()
		log.Println("Target opened a TCP connection to not request file. Closing it.")
		return
	}

	key := messages.RequestFile(msg).Key()
	file, found := t.fileManager.Find(key)

	if !found {
		conn.Write(messages.NewFileNotFound(t.addr.IP, key))
		conn.Close()
		return
	}

	_, err = conn.Write(messages.NewGetFile(t.addr.IP, key))
	if err != nil {
		log.Println("Error trying to tranfer file. Close the connection")
		conn.Close()
		return
	}

	file.Open()
	s := &stream{
		file:       file,
		bufferSize: DEFAULT_BUFFER_SIZE,
		stopFlag:   false,
		conn:       conn,
	}

	t.uploads[key] = s
	go s.upload()
}
