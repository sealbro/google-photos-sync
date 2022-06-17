package model

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"time"
)

type User struct {
	Id        string
	Created   time.Time
	Modified  time.Time
	FromToken []byte
	ToToken   []byte
}

func (u *User) Token(accountType AccountType) *oauth2.Token {
	token := oauth2.Token{}
	if accountType == From {
		json.Unmarshal(u.FromToken, &token)
	} else {
		json.Unmarshal(u.ToToken, &token)
	}

	return &token
}

func (u *User) CanSync() bool {
	return len(u.FromToken) > 0 && len(u.ToToken) > 0
}
