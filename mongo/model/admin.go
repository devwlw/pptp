package model

import (
	"mail/admin/tool"
)

type Admin struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
	PassWord string `json:"password"`
	Avatar   string `json:"avatar"`
}

func (a Admin) CheckPassword(password, cipher string) bool {
	return tool.Md5(password) == cipher
}
