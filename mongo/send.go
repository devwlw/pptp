package mongo

import (
	"context"
	"fmt"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

//Get(userName string) (*Admin, error)
//Update(*Admin) error

type Send struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewSend(collection string, db *mongo.Database) *Send {
	ctx := context.Background()
	return &Send{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (a *Send) Get(start, limit int) ([]*model.SendInfo, int64, error) {
	option := new(options.FindOptions)
	option.SetLimit(int64(limit))
	option.SetSkip(int64(start))
	cur, err := a.c.Find(a.ctx, bson.M{}, option)
	if err != nil {
		return nil, 0, err
	}
	re := make([]*model.SendInfo, 0)
	for cur.Next(a.ctx) {
		sender := new(model.SendInfo)
		err := cur.Decode(sender)
		if err != nil {
			return nil, 0, err
		}
		re = append(re, sender)
	}
	total, err := a.c.CountDocuments(a.ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return re, total, nil
}

func (a *Send) GetAll() ([]*model.SendInfo, error) {
	cur, err := a.c.Find(a.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	re := make([]*model.SendInfo, 0)
	for cur.Next(a.ctx) {
		sender := new(model.SendInfo)
		err := cur.Decode(sender)
		if err != nil {
			return nil, err
		}
		re = append(re, sender)
	}
	return re, nil
}

func (a *Send) DeleteById(id string) error {
	_, err := a.c.DeleteOne(a.ctx, bson.M{"id": id})
	return err
}

func (a *Send) GetById(id string) (*model.SendInfo, error) {
	sender := new(model.SendInfo)
	err := a.c.FindOne(a.ctx, bson.M{"id": id}).Decode(sender)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return sender, nil
}

func (a *Send) Upsert(sender *model.SendInfo) error {
	old, err := a.GetById(sender.Id)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = a.c.InsertOne(a.ctx, sender)
		if err != nil {
			return err
		}
	}

	return a.c.FindOneAndReplace(a.ctx, bson.M{"id": sender.Id}, sender).Err()
}

func (a *Send) Flush() error {
	re, err := a.c.DeleteMany(a.ctx, bson.D{})
	if err != nil {
		return err
	}
	fmt.Printf("delete %d sender", re.DeletedCount)
	return nil
}
