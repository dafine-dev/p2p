package transfer

import (
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
	index := int64(0)

	defer s.file.Close()
	defer syscall.Close(s.sock)

	for {
		if s.stopFlag {
			return
		}

		buffer := make([]byte, s.bufferSize)
		n, err := syscall.Read(s.sock, buffer)
		if err != nil {
			return
		}

		n, err = s.file.WriteAt(buffer[:n], index)
		if err != nil {
			return
		}

		index += int64(n)
	}
}

func (s *stream) upload() {
	retries := 0
	index := int64(0)
	defer s.file.Close()
	defer syscall.Close(s.sock)

	for {
		if s.stopFlag {
			return
		}

		buffer := make([]byte, s.bufferSize)
		n, err := s.file.ReadAt(buffer, index)
		if err != nil {
			return
		}

		n, err = syscall.Write(s.sock, buffer)
		if err != nil {
			if retries < 3 {
				retries++
				continue
			} else {
				return
			}
		}

		retries = 0
		index += int64(n)
	}
}

func (s *stream) Stop() {
	s.stopFlag = false
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
