package mongo

import (
	"context"
	"mail/mongo/model"

	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Deploy struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewDeploy(collection string, db *mongo.Database) *Deploy {
	ctx := context.Background()
	return &Deploy{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (d *Deploy) Find() (*model.Deploy, error) {
	deploy := new(model.Deploy)
	err := d.c.FindOne(d.ctx, bson.D{}).Decode(deploy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return deploy, nil
}

func (d *Deploy) Upsert(deploy *model.Deploy) error {
	old, err := d.Find()
	if err != nil {
		return err
	}
	if old == nil {
		deploy.Id = uuid.NewV4().String()
		_, err = d.c.InsertOne(d.ctx, deploy)
		if err != nil {
			return err
		}
	}
	return d.c.FindOneAndReplace(d.ctx, bson.D{}, deploy).Err()
}
