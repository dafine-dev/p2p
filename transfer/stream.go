package transfer

import (
	"fmt"
	"log"
	"p2p/files"
	"p2p/shared"
	"sync"
	"syscall"
)

type stream struct {
	bufferSize  int
	sock        shared.Socket
	file        *files.File
	stopFlag    bool
	bufferMutex sync.Mutex
}

func (s *stream) download() {
	defer log.Println("finished downloading")
	defer s.file.Close()
	defer syscall.Close(s.sock)

	log.Println("download")
	index := int64(0)
	s.file.Create()
	for {
		if s.stopFlag {
			return
		}

		buffer := make([]byte, s.bufferSize)
		n, err := syscall.Read(s.sock, buffer)
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
	defer log.Println("finished uploading")
	defer s.file.Close()
	defer syscall.Close(s.sock)

	log.Println("upload")
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

		fmt.Println(string(buffer))
		n, err = syscall.Write(s.sock, buffer)
		if err != nil {
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
