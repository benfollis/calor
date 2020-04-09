package database

import (
	"database/sql"
	"github.com/benfollis/calor/internal/thermometers"
)


// A CalorDB implements the functions required to record and read back
// readings from calor thermometers
type CalorDB interface {
	DB() *sql.DB
	Init()
	Latest(thermometer string) (thermometers.Reading, error)
	InsertReading(reading thermometers.Reading)
	Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error)
}

// Some DB drivers are weird about dates (SQLITE), hence when we query
// we'd like to do so with integers representing exact unix timestamp
type UnixTimestampRange struct {
	Begin int64
	End int64
}

