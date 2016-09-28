package models

type ClientToken struct {
	Common
	User uint `json:"user"`
	Token string `json:"token"`
	Client uint `json:"client"`
}

func (target *ClientToken) Marshal() interface{} {
	return target
}
