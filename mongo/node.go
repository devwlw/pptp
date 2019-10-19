package mongo

import (
	"context"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Node struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewNode(collection string, db *mongo.Database) *Node {
	ctx := context.Background()
	return &Node{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (d *Node) Find() (*model.Node, error) {
	node := new(model.Node)
	err := d.c.FindOne(d.ctx, bson.D{}).Decode(node)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return node, nil
}

func (d *Node) Upsert(node *model.Node) error {
	old, err := d.Find()
	if err != nil {
		return err
	}
	if old == nil {
		_, err = d.c.InsertOne(d.ctx, node)
		if err != nil {
			return err
		}
	}

	return d.c.FindOneAndReplace(d.ctx, bson.D{}, node).Err()
}
