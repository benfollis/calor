package database

import (
	"database/sql"
	"follis.net/internal/thermometers"
)

type CalorDB interface {
	Open() *sql.DB
	Init()
	Latest(thermometer string) (thermometers.Reading, error)
	InsertReading(reading thermometers.Reading)
	Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error)
}

type UnixTimestampRange struct {
	Begin int64
	End int64
}

