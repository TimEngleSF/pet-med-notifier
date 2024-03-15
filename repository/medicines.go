package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeKey struct {
	Hour int `bson:"hour"`
	Min  int `bson:"min"`
}

type Medicine struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	TimeToTake *TimeKey           `bson:"time-to-take"`
	Taken      bool               `bson:"taken"`
	TimeTaken  *time.Time         `bson:"time-taken"`
	Date       *time.Time         `bson:"date"`
	Due        bool               `bson:"overdue"`
}

type Medicines []Medicine

type GroupedMedicines map[TimeKey][]Medicine

func GetDailyMedicines(c context.Context, d mongo.Database) (Medicines, error) {
	collection := d.Collection("medicines")
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Filter by Date
	filter := bson.M{"date": bson.M{"$gte": today}}

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
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	m.Date = &today

	result, err := collection.InsertOne(c, m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t TimeKey) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Min)
}

func (gm GroupedMedicines) SortKeys() []TimeKey {
	var keys []TimeKey
	for key := range gm {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Hour < keys[j].Hour {
			return true
		}
		if keys[i].Hour > keys[j].Hour {
			return false
		}
		return keys[i].Min < keys[j].Min
	})

	return keys
}

func (meds Medicines) GroupByTime() GroupedMedicines {
	medicineGroups := GroupedMedicines{}
	for _, med := range meds {
		if med.TimeToTake != nil {
			key := *med.TimeToTake
			medicineGroups[key] = append(medicineGroups[key], med)
		}
	}

	return medicineGroups
}
