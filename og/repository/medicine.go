package repository

import (
	"context"
	db "lily-med/DB"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeOfDay struct {
	Hour int
	Min  int
}

type Medicine struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty" `
	Name       string              `bson:"name"`
	Taken      bool                `bson:"taken"`
	Disabled   bool                `bson:"disabled"`
	TimeToTake TimeOfDay           `bson:"timeToTake"`
	TimeTaken  *TimeOfDay          `bson:"timeTaken,omitempty"`
}

func FindMedicineById(ctx context.Context, _id primitive.ObjectID) (*Medicine, error) {
	d, err := db.GetInstance(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}

	coll := d.Db.Collection(d.Collections.Meds)
	res := coll.FindOne(ctx, filter)

	var medicine Medicine
	if err := res.Decode(&medicine); err != nil {
		log.Printf("Error decoding medicine in FindMedicineById: %v\n", err)
		return nil, err
	}

	return &medicine, nil
}

func (m *Medicine) AddMedicine(ctx context.Context) (*mongo.InsertOneResult, error) {

	d, err := db.GetInstance(ctx)
	if err != nil {
		return nil, err
	}
	coll := d.Db.Collection(d.Collections.Meds)
	res, err := coll.InsertOne(ctx, m)
	if err != nil {
		log.Printf("Error adding medicine: %v\n", err)
		return nil, err
	}

	return res, nil
}
