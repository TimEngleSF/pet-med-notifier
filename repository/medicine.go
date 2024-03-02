package repository

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeMed struct {
	Hour int
	Min  int
}

type Medicine struct {
	Name       string
	TimeToTake *TimeMed
	Taken      bool
	TimeTaken  *TimeMed
	Date       *time.DateOnly
}

func (m *Medicine) CreateMedicine(c echo.Context, db mongo.Database) error {}
