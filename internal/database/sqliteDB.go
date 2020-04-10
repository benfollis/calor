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
	database *sql.DB
}

func CreateSqliteDB(file string) SqliteDB{
	db, err := sql.Open("sqlite3", file)
	utils.CheckPanic(err)
	sqlite := SqliteDB{
		database: db,
	}
	return sqlite
}

func (sqlite SqliteDB) DB() *sql.DB {
	return sqlite.database
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

func (sqlite SqliteDB) Init() {
	CreateTable(sqlite.DB(), createTableAndIndexes)
}


var insertRow = `
	INSERT INTO readings (name, unit, temperature, temperature_kelvin, unixtime)
	values (?, ?, ?, ?, ?)
`

func (sqldb SqliteDB) InsertReading(reading thermometers.Reading) {
	db := sqldb.database
	var kelvin float64
	switch reading.Unit {
	case "K":
		kelvin = reading.Temp
	case "F":
		kelvin = (((reading.Temp - 32) * 5) / 9) + 273.15
	case "C":
		kelvin = reading.Temp + 273.15
	}
	tx, stmt := PrepareStatement(db, insertRow)
	_, err := stmt.Exec(reading.Name, reading.Unit, reading.Temp, kelvin, reading.Time.Unix())
	utils.CheckLog(err)
	defer stmt.Close()
	tx.Commit()
}

const sqliteLastReading = `
	SELECT name, unit, temperature, unixtime  
	FROM readings 
	WHERE name = ? 
	ORDER BY id
	DESC
	LIMIT 1
`

func sqliteMakeReading(rows *sql.Rows) thermometers.Reading {
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
	return FetchLatest(sqldb.database, sqliteLastReading, thermometer, sqliteMakeReading)
}

const sqliteReadingsBetween = `
SELECT name, unit, temperature, unixtime  
	FROM readings 
	WHERE name = ?
	AND unixtime >= ?
	AND unixtime <= ?
	ORDER BY id
`
func (sqldb SqliteDB) Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error) {
	db := sqldb.DB()
	tx, stmt := PrepareStatement(db, sqliteReadingsBetween)
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
		readings = append(readings, sqliteMakeReading(rows))
	}
	tx.Commit()
	return readings, searchError
}