package model

type GoogleCallback struct {
	Status AccountType `query:"state" json:"state"`
	Code   string      `query:"code" json:"code"`
}
