package storage

import "time"

type User struct {
	Id      int
	Name    string
	Email   string
	RegDate time.Time
}
