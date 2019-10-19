package model

//实例信息
type Container struct {
	Id   string          `json:"id"`
	Ip   string          `json:"ip"`
	List []ContainerInfo `json:"list"`
}

type ContainerInfo struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	Ip          string `json:"ip"`
	PptpIp      string `json:"pptpIp"`
	CreatedTime int64  `json:"createdTime"`
	Status      string `json:"status"`
}
