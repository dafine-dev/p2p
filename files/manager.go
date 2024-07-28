package files

import (
	"crypto/sha1"
	"math/big"
	"os"
	"p2p/shared"
	"path/filepath"
)

type Manager struct {
	directory string
	all       map[shared.HashKey]*File
}

func NewManager(path string) *Manager {
	return &Manager{
		directory: path,
		all:       make(map[shared.HashKey]*File),
	}
}

func (m *Manager) Find(key shared.HashKey) (*File, bool) {
	file, found := m.all[key]
	return file, found
}

func (m *Manager) New(name string) *File {
	fullname := filepath.Join(m.directory, name)

	hasher := sha1.New()
	hasher.Write([]byte(name))
	key := shared.HashKey(hasher.Sum(nil))

	file := &File{
		Name: fullname,
		Key:  key,
		Id:   new(big.Int).SetBytes(key[:]),
	}

	m.all[key] = file

	return file
}

func (m *Manager) Get(name string) *File {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	key := shared.HashKey(hasher.Sum(nil))

	return m.all[key]
}

func (m *Manager) SetUp() {
	entries, err := os.ReadDir(m.directory)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			m.New(entry.Name())
		}
	}
}
