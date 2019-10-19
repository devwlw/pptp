package mongo

import (
	"context"
	"mail/mongo/model"

	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Variable struct {
	Collection string
	DB         *mongo.Database
	ctx        context.Context
	c          *mongo.Collection
}

func NewVariable(collection string, db *mongo.Database) *Variable {
	ctx := context.Background()
	return &Variable{
		Collection: collection,
		DB:         db,
		ctx:        ctx,
		c:          db.Collection(collection),
	}
}

func (d *Variable) Find() (*model.Variable, error) {
	variable := new(model.Variable)
	err := d.c.FindOne(d.ctx, bson.D{}).Decode(variable)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return variable, nil
}

func (d *Variable) Upsert(variable *model.Variable) error {
	old, err := d.Find()
	if err != nil {
		return err
	}
	if old == nil {
		variable.Id = uuid.NewV4().String()
		_, err = d.c.InsertOne(d.ctx, variable)
		if err != nil {
			return err
		}
	}
	return d.c.FindOneAndReplace(d.ctx, bson.D{}, variable).Err()
}
