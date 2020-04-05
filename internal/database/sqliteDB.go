package database

import (
	"database/sql"
	"errors"
	"follis.net/internal/thermometers"
	"follis.net/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type SqliteDB struct {
	DBFile string
}

func (sqlite SqliteDB) Open() *sql.DB {
	db, err := sql.Open("sqlite3", sqlite.DBFile)
	utils.Check(err)
	return db
}

func (sqlite SqliteDB) Init() {
	db := sqlite.Open()
	defer db.Close()
	createTable(db)
}

var createTableAndIndexes = `
	CREATE TABLE IF NOT EXISTS readings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		unit TEXT NOT NULL,
		temperature REAL NOT NULL,
		temperature_kelvin REAL NOT NULL,
		unixtime INTEGER NOT NULL
	);
`

func createTable(db *sql.DB) {
	_, err := db.Exec(createTableAndIndexes)
	utils.Check(err)
}

var insertRow = `
	INSERT INTO readings (name, unit, temperature, temperature_kelvin, unixtime)
	values (?, ?, ?, ?, ?)
`

func prepareStatemt(db *sql.DB, statement string) (*sql.Tx, *sql.Stmt) {
	tx, err := db.Begin()
	utils.Check(err)
	stmt, err := tx.Prepare(statement)
	utils.Check(err)
	return tx, stmt
}

func (sqldb SqliteDB) InsertReading(reading thermometers.Reading) {
	db := sqldb.Open()
	defer db.Close()
	var kelvin float64
	switch reading.Unit {
	case "K":
		kelvin = reading.Temp
	case "F":
		kelvin = (((reading.Temp - 32) * 5) / 9) + 273.15
	case "C":
		kelvin = reading.Temp + 273.15
	}
	tx, stmt := prepareStatemt(db, insertRow)
	_, err := stmt.Exec(reading.Name, reading.Unit, reading.Temp, kelvin, reading.Time.Unix())
	utils.Check(err)
	defer stmt.Close()
	tx.Commit()
}

var selectLastReading = `
	SELECT name, unit, temperature, unixtime  
	FROM readings 
	WHERE name = ? 
	ORDER BY id
	DESC
	LIMIT 1
`

func makeReading(rows *sql.Rows) thermometers.Reading {
	var name, unit string
	var temperature float64
	var unixtime int64
	err := rows.Scan(&name, &unit, &temperature, &unixtime)
	utils.Check(err)
	result := thermometers.Reading{
		Temp: temperature,
		Unit: unit,
		Name: name,
		Time: time.Unix(unixtime, 0),
	}
	return result
}

func (sqldb SqliteDB) Latest(thermometer string) (thermometers.Reading, error) {
	db := sqldb.Open()
	defer db.Close()
	tx, stmt := prepareStatemt(db, selectLastReading)
	rows, err := stmt.Query(thermometer)
	utils.Check(err)
	defer stmt.Close()
	defer rows.Close()
	searchError := errors.New("404")
	if rows.Next() { // should have only one
		return makeReading(rows), nil
	}
	tx.Commit()
	return thermometers.Reading{}, searchError
}

var selectReadingsBetween = `
SELECT name, unit, temperature, unixtime  
	FROM readings 
	WHERE name = ?
	AND unixtime >= ?
	AND unixtime <= ?
	ORDER BY id
`
func (sqldb SqliteDB) Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error) {
	db := sqldb.Open()
	defer db.Close()
	tx, stmt := prepareStatemt(db, selectReadingsBetween)
	defer stmt.Close()
	end := timestampRange.End
	if end == 0 {
		end = int64(time.Now().Unix())
	}
	rows, err := stmt.Query(name, timestampRange.Begin, end)
	defer rows.Close()
	utils.Check(err)
	readings := make([]thermometers.Reading, 0)
	searchError := errors.New("404")
	for rows.Next() {
		searchError = nil
		readings = append(readings, makeReading(rows))
	}
	tx.Commit()
	return readings, searchError
}