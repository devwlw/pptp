package context

import (
	"encoding/json"
	"io/ioutil"
	"mail/config"
	"mail/deploy/docker"
	"mail/mongo"
	"path/filepath"
)

var ContextSingle *Context

type Context struct {
	Config      *config.DeployConfig
	MongoClient *mongo.Client
	SDK         *docker.SDK
}

func InitContext(configPath string) {
	conf := new(config.DeployConfig)
	configPath, err := filepath.Abs(configPath)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		panic(err)
	}
	mongoCli, err := mongo.NewClient(conf.MongoConfig)
	if err != nil {
		panic(err)
	}
	ctx := &Context{
		Config:      conf,
		MongoClient: mongoCli,
		SDK:         docker.NewSDK(conf.DockerApiVersion, conf),
	}
	ContextSingle = ctx
}
