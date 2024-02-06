package db

import (
	"context"
	"lily-med/helpers"
	"testing"
)

func TestConnectClient(t *testing.T) {
	uri, err := helpers.GetURIString()
	if err != nil {
		t.Errorf("Error getting URI string: %v", err)
		return
	}
	client := ConnectClient(context.TODO(), uri)
	if err := client.Ping(context.TODO(), nil); err != nil {
		t.Errorf("Client failed to connect: %v", err)
	}
	client.Disconnect(context.TODO())
}

func TestConnectDatabase(t *testing.T) {
	uri, err := helpers.GetURIString()
	if err != nil {
		t.Errorf("Error getting URI string: %v", err)
		return
	}
	expectedName := "lily-med-test"
	client := ConnectClient(context.TODO(), uri)
	d := ConnectDatabase(client, expectedName)
	dName := d.Name()
	if dName != expectedName {
		t.Errorf("Expected database name to be '%v', but got '%v'", expectedName, dName)
	}

	client.Disconnect(context.TODO())
}
