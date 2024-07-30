package users

import (
	"math"
	"p2p/shared"
)

type Table struct {
	Predecessor *User
	Successor   *User
	all         map[shared.HashId]*User
	Current     *User
}

func StartTable(addr shared.Addr) *Table {
	current := New(addr)

	all := make(map[shared.HashId]*User)

	for i := 0; i < 7; i++ {
		id := current.Id + uint(math.Pow(2, float64(i)))
		id = id % shared.MaxId

		all[id] = nil
	}

	return &Table{
		all:         all,
		Current:     current,
		Successor:   current,
		Predecessor: current,
	}
}

func (t *Table) Owns(id shared.HashId) bool {
	return shared.Distance(id, t.Successor.Id) <= shared.Distance(t.Current.Id, t.Successor.Id)
}

func (t *Table) Nearest(id shared.HashId) *User {
	u := t.Successor
	for key, user := range t.all {
		if user != nil && shared.Distance(key, id) < shared.Distance(key, u.Id) {
			u = user
		}
	}

	return u
}

func (t *Table) Update(key shared.HashId, user *User) {
	if _, in := t.all[key]; in {
		t.all[key] = user
	}
}

func (t *Table) IsPredecessor(user *User) bool {
	return shared.Distance(user.Id, t.Current.Id) <
		shared.Distance(t.Predecessor.Id, t.Current.Id)
}

func (t *Table) IsSuccessor(user *User) bool {
	return shared.Distance(t.Current.Id, user.Id) <
		shared.Distance(t.Current.Id, t.Successor.Id)
}

func (t *Table) SetPredecessor(user *User) bool {
	if t.IsPredecessor(user) {
		t.Predecessor = user
		return true
	}

	return false
}

func (t *Table) SetSuccessor(user *User) bool {
	if t.IsSuccessor(user) {
		t.Successor = user
		return true
	}

	return false
}

// func (t *Table) New(addr shared.Addr) *User {
// 	hasher := sha1.New()
// 	hasher.Write(addr.Addr[:])
// 	rawid := hasher.Sum(nil)
// 	id := new(big.Int).SetBytes(rawid)
// 	new_user := &User{
// 		Addr:  addr,
// 		Id:    id,
// 		RawId: [20]byte(rawid),
// 	}
// 	t.all[addr] = new_user
// 	return new_user
// }

func (t *Table) SetCurrent(addr shared.Addr, name string) *User {
	current := New(addr)
	t.Current = current
	return current
}

// func (t *Table) Add(user *User) {
// 	t.all[] = user
// }

func (t *Table) All() map[shared.HashId]*User {
	return t.all
}
