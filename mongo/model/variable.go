package model

type Variable struct {
	Id   string            `json:"id"`
	Data map[string]string `json:"data"`
}
