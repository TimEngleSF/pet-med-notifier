package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const TimeZone = "America/Los_Angeles"

type TimeMed struct {
	Hour int `bson:"hour"`
	Min  int `bson:"min"`
}

type Medicine struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	TimeToTake *TimeMed           `bson:"time-to-take"`
	Taken      bool               `bson:"taken"`
	TimeTaken  *TimeMed           `bson:"time-taken"`
	Date       *time.Time         `bson:"date"`
}

// func (m *Medicine) CreateMedicine(c echo.Context, db mongo.Database) error {}

func GetDailyMedicines(c context.Context, d mongo.Database) ([]Medicine, error) {
	collection := d.Collection("medicines")

	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(loc).Truncate(24 * time.Hour)
	// Filter by Date
	filter := bson.M{"date": bson.M{"$gt": today}}

	var results []Medicine
	cursor, err := collection.Find(c, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var medicine Medicine
		if err := cursor.Decode(&medicine); err != nil {
			return nil, err
		}
		results = append(results, medicine)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func AddDailyMedicine(c context.Context, d *mongo.Database, m Medicine) (*mongo.InsertOneResult, error) {
	collection := d.Collection("medicines")

	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(loc).Truncate(24 * time.Hour)
	m.Date = &today
	//	m.Id = bson.TypeObjectID

	result, err := collection.InsertOne(c, m)
	if err != nil {
		return nil, err
	}
	return result, nil
}
