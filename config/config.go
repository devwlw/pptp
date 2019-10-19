package config

type AdminConfig struct {
	DefaultMaxInstance int         `json:"defaultMaxInstance"` //每台机器最大docker实例数
	DefaultMaxMachine  int         `json:"defaultMaxMachine"`  //最多的机器数
	MongoConfig        MongoConfig `json:"mongoConfig"`
	TemplatePath       string      `json:"templatePath"`
}

type MongoConfig struct {
	EndPoint            string `json:"endpoint"`
	Port                string `json:"port"`
	User                string `json:"user"`
	Password            string `json:"password"`
	DataBase            string `json:"database"`
	DeployCollection    string `json:"deployCollection"`
	NodeCollection      string `json:"nodeCollection"`
	ContainerCollection string `json:"containerCollection"`
	AdminCollection     string `json:"adminCollection"`
	TemplatesCollection string `json:"templatesCollection"`
	ReceiverCollection  string `json:"receiverCollection"`
	LogCollection       string `json:"logCollection"`
	SendCollection      string `json:"sendCollection"`
	VariableCollection  string `json:"variableCollection"`
}

type DeployConfig struct {
	DockerRepository string      `json:"dockerRepository"`
	DockerApiVersion string      `json:"dockerApiVersion"`
	AdminServerAddr  string      `json:"adminServerAddr"`
	HostIp           string      `json:"hostIp"`
	MongoConfig      MongoConfig `json:"mongoConfig"`
	PPTP             PPTP        `json:"pptp"`
}

type PPTP struct {
	Endpoint string `json:"endpoint"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
