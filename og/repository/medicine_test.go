package repository

import (
	"context"
	"fmt"
	db "lily-med/DB"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetMedicine(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("GO_ENV", "test")
	testCtx := context.Background()
	d, err := db.GetInstance(testCtx)
	if err != nil {
		t.Errorf("error getting db instance: %v\n", err)
	}
	defer d.CloseConnection() // Ensure the DB connection is closed when the test finishes

	coll := d.Db.Collection(d.Collections.Meds)

	dummyMedicine := []Medicine{
		{
			Name:       "TestMedicine",
			Taken:      false,
			Disabled:   false,
			TimeToTake: TimeOfDay{Hour: 10, Min: 00},
		},
		{
			Name:       "TestMedicine2",
			Taken:      true,
			Disabled:   false,
			TimeToTake: TimeOfDay{Hour: 11, Min: 30},
		},
	}

	for _, d := range dummyMedicine {
		insertRes, err := coll.InsertOne(testCtx, d)
		if err != nil {
			t.Fatalf("failed to insert dummyMedicine: %v", err)
		}

		idStr := fmt.Sprintf("%v", insertRes.InsertedID)
		fmt.Println(idStr)

		objectId, ok := insertRes.InsertedID.(primitive.ObjectID)
		if !ok {
			t.Errorf("insertedId failed type check\n")
		}

		med, err := FindMedicineById(testCtx, objectId)
		if err != nil {
			t.Fatalf("FindMedicineById failed: %v", err)
		}

		if med.Name != d.Name {
			t.Errorf("expected Name to be %s but got %s", d.Name, med.Name)
		}

		// Ensure cleanup after each test case
		if _, err := coll.DeleteOne(testCtx, bson.M{"_id": objectId}); err != nil {
			t.Logf("failed to clean up test medicine: %v", err)
		}
	}
}
