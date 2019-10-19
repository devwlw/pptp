package model

//已注册的节点信息
type NodeInfo struct {
	Ip      string `json:"ip"`
	Aux     string `json:"aux"`
	IpRange string `json:"ipRange"`
	Subnet  string `json:"subnet"`
	Nic     string `json:"nic"`
	//HostName    string `json:"hostname"`
	CreatedTime int64 `json:"createdTime"`
}

type Node struct {
	Nodes []NodeInfo `json:"nodes"`
}
