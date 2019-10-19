package context

import (
	"encoding/json"
	"io/ioutil"
	"mail/config"
	"mail/mongo"
	"path/filepath"
)

var ContextSingle *Context

type Context struct {
	Config      *config.AdminConfig
	MongoClient *mongo.Client
}

func InitContext(configPath string) {
	conf := new(config.AdminConfig)
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
	}
	ContextSingle = ctx
}
