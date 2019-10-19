package model

//docker基本部署信息
type Deploy struct {
	Id          string `json:"id"`
	MaxInstance int    `json:"maxInstance"` //每台机器最大docker实例数
	MaxMachine  int    `json:"maxMachine"`  //最多的机器数
	RefreshTick int    `json:"refreshTick"` //后台获取docker信息的更新时间
}
