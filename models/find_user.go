package models

type FindUser struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}
