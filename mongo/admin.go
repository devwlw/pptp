package mongo

import (
	"context"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Get(userName string) (*Admin, error)
//Update(*Admin) error

type Admin struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewAdmin(collection string, db *mongo.Database) *Admin {
	ctx := context.Background()
	return &Admin{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (a *Admin) Get(userName string) (*model.Admin, error) {
	admin := new(model.Admin)
	err := a.c.FindOne(a.ctx, bson.M{"userName": userName}).Decode(admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return admin, nil
}

func (a *Admin) Update(admin *model.Admin) error {
	old, err := a.Get(admin.UserName)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = a.c.InsertOne(a.ctx, admin)
		if err != nil {
			return err
		}
	}
	return a.c.FindOneAndReplace(a.ctx, bson.M{"userName": admin.UserName}, admin).Err()
}
