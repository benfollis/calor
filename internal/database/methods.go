package database

import (
	"database/sql"
	"calor/internal/thermometers"
)

type CalorDB interface {
	DB() *sql.DB
	Init()
	Latest(thermometer string) (thermometers.Reading, error)
	InsertReading(reading thermometers.Reading)
	Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error)
}

type UnixTimestampRange struct {
	Begin int64
	End int64
}

