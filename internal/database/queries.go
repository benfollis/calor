package database

import (
	"database/sql"
	"follis.net/internal/thermometers"
	_ "github.com/mattn/go-sqlite3"
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
	tx, err := db.Begin()
	check(err)
	stmt, err := tx.Prepare(insertRow)
	check(err)
	stmt.Exec(reading.Name, reading.Unit, reading.Temp, kelvin, reading.Time)
	defer stmt.Close()
	tx.Commit()
}