package files

import (
	"crypto/sha1"
	"os"
)

type Hash = [160]byte

type File struct {
	os.File
	Name   string
	Key    Hash
	Status int
}

var allFiles map[Hash]*File

func SetUp() {

}

func Search(key Hash) (*File, bool) {
	file, found := allFiles[key]
	return file, found
}

func New(name string) *File {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	hasher := sha1.New()
	hasher.Write([]byte(name))

	file := &File{
		File: *f,
		Name: name,
		Key:  Hash(hasher.Sum(nil)),
	}

	return file
}

func All() []*File {
	all := make([]*File, 0)
	for _, value := range allFiles {
		all = append(all, value)
	}

	return all
}
