package database

import (
	"database/sql"
	"errors"
	"github.com/benfollis/calor/internal/thermometers"
	"github.com/benfollis/calor/internal/utils"
)


// A CalorDB implements the functions required to record and read back
// readings from calor thermometers
type CalorDB interface {
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


// CreateTable creates the table(s) required for calor
func CreateTable(db *sql.DB, createString string) {
	_, err := db.Exec(createString)
	utils.CheckLog(err)
}


// Prepare statements prepares the given statement for a transaction
func PrepareStatement(db *sql.DB, statement string) (*sql.Tx, *sql.Stmt) {
	tx, err := db.Begin()
	utils.CheckLog(err)
	stmt, err := tx.Prepare(statement)
	utils.CheckLog(err)
	return tx, stmt
}

// Fetches a single reading from the DB
func FetchLatest(db *sql.DB, querySQL string, thermometer string, readingMaker func (rows *sql.Rows) thermometers.Reading) (thermometers.Reading, error) {
	tx, stmt := PrepareStatement(db, querySQL)
	rows, err := stmt.Query(thermometer)
	utils.CheckLog(err)
	defer stmt.Close()
	defer rows.Close()
	searchError := errors.New("404")
	if rows.Next() { // should have only one
		return readingMaker(rows), nil
	}
	tx.Commit()
	return thermometers.Reading{}, searchError
}
