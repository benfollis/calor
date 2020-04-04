package database

import (
	"database/sql"
	"follis.net/internal/thermometers"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}


var createTableAndIndexes = `
	CREATE TABLE IF NOT EXISTS readings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		unit TEXT NOT NULL,
		temperature REAL NOT NULL,
		temperature_kelvin REAL NOT NULL,
		timestamp INTEGER NOT NULL
	);
`

func CreateTable(db *sql.DB) {
	_, err := db.Exec(createTableAndIndexes)
	check(err)
}

var insertRow = `
	INSERT INTO readings (name, unit, temperature, temperature_kelvin, timestamp)
	values (?, ?, ?, ?, ?)
`

func prepareStatemt(db *sql.DB, statement string) (*sql.Tx, *sql.Stmt) {
	tx, err := db.Begin()
	check(err)
	stmt, err := tx.Prepare(statement)
	check(err)
	return tx, stmt
}

func InsertReading(db *sql.DB, reading thermometers.Reading) {
	var kelvin float64
	switch reading.Unit {
	case "K":
		kelvin = reading.Temp
	case "F":
		kelvin = (((reading.Temp -32) * 5) / 9) + 273.15
	case "C":
		kelvin = reading.Temp + 273.15
	}
	tx, stmt := prepareStatemt(db, insertRow)
	stmt.Exec(reading.Name, reading.Unit, reading.Temp, kelvin, reading.Time)
	defer stmt.Close()
	tx.Commit()
}

var selectLastReading = `
	SELECT (name, unit, temperature, timestamp)  
	FROM readings 
	WHERE name = ? 
	ORDER BY id
	DESC
	LIMIT 1
`

func LastReading(db *sql.DB, name string) thermometers.Reading {
   tx, stmt := prepareStatemt(db, selectLastReading)
   rows, err := stmt.Query(name)
   var result thermometers.Reading
   check(err)
   defer stmt.Close()
   defer rows.Close()
   var thermometer, unit string
   var temperature float64
   var timestamp time.Time
   for rows.Next() {  // should have only one
   	err := rows.Scan(&thermometer, &unit, &temperature, &timestamp)
   	check(err)
   	result = thermometers.Reading{
		Temp: temperature,
		Unit: unit,
		Name: thermometer,
		Time: timestamp,
	}
   }
   tx.Commit()
   return result
}