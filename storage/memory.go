package storage

import (
	"errors"
	"sync"
	"time"
)

type User struct {
	sync.Mutex
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"`
	RegDate string `json:"reg_date"`
}

func (u *User) Add(users map[string]User) error {
	u.Lock()
	if _, ok := users[u.Email]; ok {
		return errors.New("Email " + u.Email + " already exists")
	}

	users[u.Email] = User{
		Id:      len(users) + 1,
		RegDate: time.Now().Format("2006-02-01 15:04"),
		Name:    u.Name,
	}
	u.Unlock()

	return nil
}
