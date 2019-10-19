package model

type ReceiverInfo struct {
	Id     string `json:"id"`
	Email  string `json:"email"`
	IsUsed bool   `json:"isUsed"`
}
