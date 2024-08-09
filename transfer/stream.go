package transfer

import (
	"log"
	"net"
	"p2p/files"
	"sync"
)

type stream struct {
	bufferSize  int
	conn        *net.TCPConn
	file        *files.File
	stopFlag    bool
	bufferMutex sync.Mutex
}

func (s *stream) download() {
	defer s.file.Close()
	defer s.conn.Close()

	index := int64(0)
	s.file.Create()
	for {
		if s.stopFlag {
			return
		}

		buffer := make([]byte, s.bufferSize)
		n, err := s.conn.Read(buffer)
		if err != nil {
			s.Stop()
		}

		n, err = s.file.WriteAt(buffer[:n], index)
		if err != nil {
			log.Println("download write", err)
		}

		index += int64(n)
	}
}

func (s *stream) upload() {
	defer s.file.Close()
	defer s.conn.Close()

	// retries := 0
	index := int64(0)
	s.file.Open()

	for {
		if s.stopFlag {
			return
		}

		buffer := make([]byte, s.bufferSize)
		n, err := s.file.ReadAt(buffer, index)
		if err != nil {
			s.Stop()
		}

		n, err = s.conn.Write(buffer[:n])
		if err != nil {
			log.Println("Failed to write to UDP socket.")
			return
		}

		// retries = 0
		index += int64(n)
	}
}

func (s *stream) Stop() {
	s.stopFlag = true
}

func (s *stream) IncreaseSpeed(times int) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()
	if times < 1 {
		s.bufferSize = DEFAULT_BUFFER_SIZE
	} else {
		s.bufferSize = times * DEFAULT_BUFFER_SIZE
	}
}
