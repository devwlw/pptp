package mongo

import (
	"context"
	"fmt"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Get(userName string) (*Admin, error)
//Update(*Admin) error

type Receiver struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewReceiver(collection string, db *mongo.Database) *Receiver {
	ctx := context.Background()
	return &Receiver{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (r *Receiver) Get(start, limit int) ([]*model.ReceiverInfo, int64, error) {
	option := new(options.FindOptions)
	option.SetLimit(int64(limit))
	option.SetSkip(int64(start))
	cur, err := r.c.Find(r.ctx, bson.M{}, option)
	if err != nil {
		return nil, 0, err
	}
	re := make([]*model.ReceiverInfo, 0)
	for cur.Next(r.ctx) {
		receiver := new(model.ReceiverInfo)
		err := cur.Decode(receiver)
		if err != nil {
			return nil, 0, err
		}
		re = append(re, receiver)
	}
	total, err := r.c.CountDocuments(r.ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return re, total, nil
}

func (r *Receiver) GetAll() ([]*model.ReceiverInfo, error) {
	cur, err := r.c.Find(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	re := make([]*model.ReceiverInfo, 0)
	for cur.Next(r.ctx) {
		receiver := new(model.ReceiverInfo)
		err := cur.Decode(receiver)
		if err != nil {
			return nil, err
		}
		re = append(re, receiver)
	}
	return re, nil
}

func (r *Receiver) DeleteById(id string) error {
	_, err := r.c.DeleteOne(r.ctx, bson.M{"id": id})
	return err
}

func (r *Receiver) GetById(id string) (*model.ReceiverInfo, error) {
	receiver := new(model.ReceiverInfo)
	err := r.c.FindOne(r.ctx, bson.M{"id": id}).Decode(receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return receiver, nil
}

func (r *Receiver) Upsert(receiver *model.ReceiverInfo) error {
	old, err := r.GetById(receiver.Id)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = r.c.InsertOne(r.ctx, receiver)
		if err != nil {
			return err
		}
	}

	return r.c.FindOneAndReplace(r.ctx, bson.M{"id": receiver.Id}, receiver).Err()
}

func (r *Receiver) Flush() error {
	re, err := r.c.DeleteMany(r.ctx, bson.D{})
	if err != nil {
		return err
	}
	fmt.Printf("delete %d receiver", re.DeletedCount)
	return nil
}
