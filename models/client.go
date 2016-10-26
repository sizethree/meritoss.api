package models

import "fmt"

type Client struct {
	Common
	Name         string `json:"name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"-"`
	RedirectUri  string `json:"redirect_uri"`
}

func (client Client) Url() string {
	return fmt.Sprintf("/clients?filter[id]=eq(%d)", client.ID)
}

func (client Client) Type() string {
	return "application/vnd.miritos.client+json"
}
