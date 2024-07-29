package users

import (
	"crypto/sha1"
	"math/big"
	"p2p/shared"
)

type Table struct {
	predecessor *User
	Successor   *User
	all         map[shared.Addr]*User
	Current     *User
}

func StartTable(addr shared.Addr) *Table {
	hasher := sha1.New()
	hasher.Write(addr.Addr[:])
	raw := hasher.Sum(nil)
	current := &User{
		Addr:  addr,
		Id:    new(big.Int).SetBytes(raw),
		RawId: [20]byte(raw),
	}

	return &Table{
		all:       make(map[shared.Addr]*User),
		Current:   current,
		Successor: current,
	}
}

func (t *Table) Owns(id shared.HashId) bool {
	a := shared.Distance(t.Current.Id, t.Successor.Id)
	b := shared.Distance(id, t.Successor.Id)
	return a.Cmp(b) < 0
}

func (t *Table) Nearest(key shared.HashId) *User {
	minValue := shared.MaxId
	var u *User

	for _, user := range t.all {
		if distance := shared.Distance(user.Id, key); distance.Cmp(minValue) <= 0 {
			u = user
			minValue = distance
		}
	}

	return u
}

func (t *Table) IsPredecessor(user *User) bool {
	newDistance := shared.Distance(user.Id, t.Current.Id)
	currentDistance := shared.Distance(t.predecessor.Id, t.Current.Id)
	return newDistance.Cmp(currentDistance) < 0
}

func (t *Table) IsSuccessor(user *User) bool {
	if t.Successor == t.Current {
		return true
	}
	newDistance := shared.Distance(t.Current.Id, user.Id)
	currentDistance := shared.Distance(t.Current.Id, t.Successor.Id)
	return newDistance.Cmp(currentDistance) < 0
}

func (t *Table) SetPredecessor(user *User) bool {
	if t.IsPredecessor(user) {
		t.predecessor = user
		return true
	}

	return false
}

func (t *Table) SetSuccessor(user *User) bool {
	if t.IsSuccessor(user) {
		t.Successor = user
		t.Add(user)
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

func (t *Table) Add(user *User) {
	t.all[user.Addr] = user
}
