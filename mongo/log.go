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

type Log struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewLog(collection string, db *mongo.Database) *Log {
	ctx := context.Background()
	return &Log{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (a *Log) GetById(id string) (*model.Log, error) {
	log := new(model.Log)
	err := a.c.FindOne(a.ctx, bson.M{"id": id}).Decode(log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return log, nil
}

func (a *Log) Get(start, limit int) ([]*model.Log, int64, error) {
	option := new(options.FindOptions)
	option.SetLimit(int64(limit))
	option.SetSkip(int64(start))
	cur, err := a.c.Find(a.ctx, bson.M{}, option)
	if err != nil {
		return nil, 0, err
	}
	re := make([]*model.Log, 0)
	for cur.Next(a.ctx) {
		l := new(model.Log)
		err := cur.Decode(l)
		if err != nil {
			return nil, 0, err
		}
		re = append(re, l)
	}
	total, err := a.c.CountDocuments(a.ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return re, total, nil
}

func (a *Log) Upsert(log *model.Log) error {
	old, err := a.GetById(log.Id)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = a.c.InsertOne(a.ctx, log)
		if err != nil {
			return err
		}
	}

	return a.c.FindOneAndReplace(a.ctx, bson.M{"id": log.Id}, log).Err()
}

func (a *Log) DeleteById(id string) error {
	_, err := a.c.DeleteOne(a.ctx, bson.M{"id": id})
	return err
}

func (a *Log) Flush() error {
	re, err := a.c.DeleteMany(a.ctx, bson.D{})
	if err != nil {
		return err
	}
	fmt.Printf("delete %d log", re.DeletedCount)
	return nil
}
