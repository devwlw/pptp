package mongo

import (
	"context"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewContainer(collection string, db *mongo.Database) *Container {
	ctx := context.Background()
	return &Container{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (d *Container) Find() ([]*model.Container, error) {
	containers := make([]*model.Container, 0)
	re, err := d.c.Find(d.ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	for re.Next(d.ctx) {
		container := new(model.Container)
		err := re.Decode(container)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}
	return containers, nil
}

func (d *Container) FindByField(key, value string) (*model.Container, error) {
	container := new(model.Container)
	err := d.c.FindOne(d.ctx, bson.M{key: value}).Decode(container)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return container, nil
}

func (d *Container) UpsertByField(key, value string, container *model.Container) error {
	old, err := d.FindByField(key, value)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = d.c.InsertOne(d.ctx, container)
		if err != nil {
			return err
		}
	}
	return d.c.FindOneAndReplace(d.ctx, bson.M{key: value}, container).Err()
}
