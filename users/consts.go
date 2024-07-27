package users

import "p2p/shared"

var user_table map[[4]byte]*user
var CURRENT_USERNAME string = "lucas =)"

func New(addr shared.Addr, id UserID, name string) *user {
	u := &user{
		Addr: addr,
		Id:   id,
		Name: name,
	}
	user_table[addr.Addr] = u
	return u
}

func Get(addr [4]byte) *user {
	return user_table[addr]
}

func All() []*user {
	users := make([]*user, 0)
	for _, u := range user_table {
		users = append(users, u)
	}

	return users
}
