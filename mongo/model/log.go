package model

type Log struct {
	Id          string        `json:"id"`
	Sender      *SendInfo     `json:"sender"`
	Receiver    *ReceiverInfo `json:"receiver"`
	Body        string        `json:"body"`
	Title       string        `json:"title"`
	MailType    string        `json:"mailType"`
	Proxy       string        `json:"proxy"`
	MachineIp   string        `json:"machineIp"`
	ContainerId string        `json:"containerId"`
	Success     bool          `json:"success"`
	IsProcess   bool          `json:"isProcess"`
	CreatedTime int64         `json:"createdTime"`
}
