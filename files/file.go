package files

import (
	"os"
	"p2p/shared"
)

const (
	SEARCHING uint8 = iota
	DOWNLOADING
	UPLOADING
	FINISHED
	NOT_FOUND
)

type File struct {
	os.File
	Name   string
	Key    shared.HashKey
	Status uint8
	Id     shared.HashId
}

func (f *File) Create() {
	conn, err := os.Create(f.Name)
	if err != nil {
		panic(err)
	}

	f.File = *conn
}

func (f *File) Open() {
	conn, err := os.Open(f.Name)
	if err != nil {
		panic(err)
	}

	f.File = *conn
}
