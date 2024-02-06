package db

import (
	"context"
	"os"
	"testing"
)

var testURI = "mongodb://localhost:27017"

func TestConnectClient(t *testing.T) {
	client, err := ConnectClient(context.TODO(), testURI)
	if err != nil {
		t.Errorf("Failed to connect to client: %v\n", err)
	}
	defer client.Disconnect(context.TODO())

	if err := client.Ping(context.TODO(), nil); err != nil {
		t.Fatalf("Client failed to connect: %v", err)
	}
}

func TestConnectDatabase(t *testing.T) {
	expectedName := "lily-med-test"
	client, err := ConnectClient(context.TODO(), testURI)
	if err != nil {
		t.Errorf("Failed to connect to client: %v\n", err)
	}
	defer client.Disconnect(context.TODO())

	d := ConnectDatabase(client, expectedName)
	dName := d.Name()
	if dName != expectedName {
		t.Errorf("Expected database name to be '%v', but got '%v'", expectedName, dName)
	}
}

func TestInitDatabase(t *testing.T) {
	testCases := []struct {
		name           string
		setEnv         bool
		envValue       string
		isTestEnv      bool
		expectedDbName string
		expectError    bool
	}{
		{
			name:           "Valid Test Environment",
			setEnv:         true,
			envValue:       "mongodb://localhost:27017",
			isTestEnv:      true,
			expectedDbName: "lily-med-test",
			expectError:    false,
		},
		{
			name:           "Valid Production Environment",
			setEnv:         true,
			envValue:       "mongodb://localhost:27017",
			isTestEnv:      false,
			expectedDbName: "lily-med",
			expectError:    false,
		},
		{
			name:        "Missing Environment Variable",
			setEnv:      false,
			expectError: true,
		},
		{
			name:        "Invalid URI",
			setEnv:      true,
			envValue:    "invalid-uri",
			isTestEnv:   false,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				os.Setenv("MONGO_URI", tc.envValue)
				defer os.Unsetenv("MONGO_URI")
			}

			d, err := initDatabase(context.Background(), tc.isTestEnv)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error but didn't receive one")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				defer d.CloseConnection()

				dbName := d.Db.Name()
				if dbName != tc.expectedDbName {
					t.Errorf("Expected db to be named '%v', got '%v'", tc.expectedDbName, dbName)
				}
			}
		})
	}
}
