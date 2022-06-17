package model

import "time"

type StateHistory struct {
	Id      int64
	Created time.Time
	UserId  string
	From    int
	To      int
}
