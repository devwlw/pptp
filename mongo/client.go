package mongo

import (
	"context"
	"fmt"
	"mail/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	conf     config.MongoConfig
	mongoCli *mongo.Client
	db       *mongo.Database
}

func NewClient(conf config.MongoConfig) (*Client, error) {
	//mongodb://user:password@localhost:27017/?authSource=admin
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	mongoCli, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=%s", conf.User, conf.Password, conf.EndPoint, conf.Port, conf.DataBase)))
	if err != nil {
		return nil, err
	}
	db := mongoCli.Database(conf.DataBase)

	cli := &Client{
		conf:     conf,
		mongoCli: mongoCli,
		db:       db,
	}
	return cli, nil
}

func (c *Client) DeployRepository() *Deploy {
	return NewDeploy(c.conf.DeployCollection, c.db)
}

func (c *Client) NodeRepository() *Node {
	return NewNode(c.conf.NodeCollection, c.db)
}

func (c *Client) VariableRepository() *Variable {
	return NewVariable(c.conf.VariableCollection, c.db)
}

func (c *Client) ContainerRepository() *Container {
	return NewContainer(c.conf.ContainerCollection, c.db)
}

func (c *Client) AdminRepository() *Admin {
	return NewAdmin(c.conf.AdminCollection, c.db)
}

func (c *Client) TemplatesRepository() *Templates {
	return NewTemplates(c.conf.TemplatesCollection, c.db)
}

func (c *Client) ReceiverRepository() *Receiver {
	return NewReceiver(c.conf.ReceiverCollection, c.db)
}

func (c *Client) SenderRepository() *Send {
	return NewSend(c.conf.SendCollection, c.db)
}

func (c *Client) LogRepository() *Log {
	return NewLog(c.conf.LogCollection, c.db)
}
