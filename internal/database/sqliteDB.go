package database

import (
	"database/sql"
	"errors"
	"github.com/benfollis/calor/internal/thermometers"
	"github.com/benfollis/calor/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// A SqliteDB is a CalorDB that is backed by a SqliteDB file
type SqliteDB struct {
	DBFile string
	database *sql.DB
}

func CreateSqliteDB(file string) SqliteDB{
	db, err := sql.Open("sqlite3", file)
	utils.CheckPanic(err)
	sqlite := SqliteDB{
		DBFile:   file,
		database: db,
	}
	return sqlite
}

func (sqlite SqliteDB) DB() *sql.DB {
	return sqlite.database
}

func (sqlite SqliteDB) Init() {
	createTable(sqlite.DB())
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
	utils.CheckLog(err)
}

var insertRow = `
	INSERT INTO readings (name, unit, temperature, temperature_kelvin, unixtime)
	values (?, ?, ?, ?, ?)
`

func prepareStatemt(db *sql.DB, statement string) (*sql.Tx, *sql.Stmt) {
	tx, err := db.Begin()
	utils.CheckLog(err)
	stmt, err := tx.Prepare(statement)
	utils.CheckLog(err)
	return tx, stmt
}

func (sqldb SqliteDB) InsertReading(reading thermometers.Reading) {
	db := sqldb.DB()
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
	utils.CheckLog(err)
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
	utils.CheckLog(err)
	result := thermometers.Reading{
		Temp: temperature,
		Unit: unit,
		Name: name,
		Time: time.Unix(unixtime, 0),
	}
	return result
}

func (sqldb SqliteDB) Latest(thermometer string) (thermometers.Reading, error) {
	db := sqldb.DB()
	tx, stmt := prepareStatemt(db, selectLastReading)
	rows, err := stmt.Query(thermometer)
	utils.CheckLog(err)
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
	db := sqldb.DB()
	tx, stmt := prepareStatemt(db, selectReadingsBetween)
	defer stmt.Close()
	end := timestampRange.End
	if end == 0 {
		end = int64(time.Now().Unix())
	}
	rows, err := stmt.Query(name, timestampRange.Begin, end)
	defer rows.Close()
	utils.CheckLog(err)
	readings := make([]thermometers.Reading, 0)
	searchError := errors.New("404")
	for rows.Next() {
		searchError = nil
		readings = append(readings, makeReading(rows))
	}
	tx.Commit()
	return readings, searchError
}