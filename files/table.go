package files

import (
	"p2p/shared"
)

type Location struct {
	Key  shared.HashKey
	Id   shared.HashId
	Addr shared.Addr
}

func NewLocation(key shared.HashKey, addr shared.Addr) *Location {
	return &Location{
		Key:  key,
		Id:   uint(key),
		Addr: addr,
	}
}

type Table struct {
	locations map[shared.HashKey]*Location
}

func NewTable() *Table {
	return &Table{
		locations: make(map[shared.HashKey]*Location),
	}
}

func (t *Table) Between(start shared.HashId, end shared.HashId) []*Location {
	locs := make([]*Location, 0)
	window := shared.Distance(start, end)

	for _, loc := range t.locations {

		if shared.Distance(loc.Id, end) < window {
			locs = append(locs, loc)
		}
	}

	return locs
}

func (t *Table) RemoveBetween(start shared.HashId, end shared.HashId) {
	window := shared.Distance(start, end)

	for key, loc := range t.locations {
		if shared.Distance(loc.Id, end) <= window {
			delete(t.locations, key)
		}
	}
}

func (t *Table) Add(locs ...*Location) {
	for _, loc := range locs {
		t.locations[loc.Key] = loc
	}
}

func (t *Table) Find(key shared.HashKey) (*Location, bool) {
	loc, ok := t.locations[key]
	return loc, ok
}
func (t *Table) All() map[shared.HashKey]*Location {

	return t.locations
}
