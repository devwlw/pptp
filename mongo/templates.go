package mongo

import (
	"context"
	"mail/mongo/model"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

//Get(userName string) (*Admin, error)
//Update(*Admin) error

type Templates struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewTemplates(collection string, db *mongo.Database) *Templates {
	ctx := context.Background()
	return &Templates{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (a *Templates) GetByName(name string) (*model.Templates, error) {
	templates := new(model.Templates)
	err := a.c.FindOne(a.ctx, bson.M{"name": name}).Decode(templates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return templates, nil
}

func (a *Templates) GetById(id string) (*model.Templates, error) {
	templates := new(model.Templates)
	err := a.c.FindOne(a.ctx, bson.M{"id": id}).Decode(templates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return templates, nil
}

func (a *Templates) Get() ([]*model.Templates, error) {
	templates := make([]*model.Templates, 0)
	re, err := a.c.Find(a.ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	for re.Next(a.ctx) {
		t := new(model.Templates)
		err := re.Decode(t)
		if err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (a *Templates) Upsert(templates *model.Templates) error {
	old, err := a.GetByName(templates.Name)
	if err != nil {
		return err
	}
	if old == nil {
		_, err = a.c.InsertOne(a.ctx, templates)
		if err != nil {
			return err
		}
	}

	return a.c.FindOneAndReplace(a.ctx, bson.M{"name": templates.Name}, templates).Err()
}

func (a *Templates) DeleteById(id string) error {
	_, err := a.c.DeleteOne(a.ctx, bson.M{"id": id})
	return err
}
