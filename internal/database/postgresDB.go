package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/benfollis/calor/internal/thermometers"
	"github.com/benfollis/calor/internal/utils"
	_ "github.com/lib/pq"
	"time"
)

type PostgresDB struct {
	database *sql.DB
}

const calorDBName = "calor"

func CreatePostgresDB(user string, password string, host string) PostgresDB {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, calorDBName)
	db, err := sql.Open("postgres", connString)
	utils.CheckPanic(err)
	psql := PostgresDB{
		database: db,
	}
	return psql
}

const postgresCreate = `
	CREATE TABLE IF NOT EXISTS readings (
		id serial PRIMARY KEY,
		name TEXT,
		temperature REAL,
		unit VARCHAR(4),
		date TIMESTAMPTZ
	);
	CREATE INDEX IF NOT EXISTS name_index ON readings (name);
	CREATE INDEX IF NOT EXISTS date_index ON readings (date);
`

func (psql PostgresDB) Init() {
	CreateTable(psql.database, postgresCreate)
}

const postgresInsertReading = `
INSERT INTO readings(name, temperature, unit, date) VALUES ($1, $2, $3, $4)
`

func (psql PostgresDB) InsertReading(reading thermometers.Reading) {
	db := psql.database
	tx, stmt := PrepareStatement(db, postgresInsertReading)
	defer stmt.Close()
	stmt.Exec(reading.Name, reading.Temp, reading.Unit, reading.Time)
	tx.Commit()
}

const postgresLatest = `
	SELECT name, temperature, unit, date
	FROM readings
	WHERE name = $1
	ORDER BY id DESC
	LIMIT 1
`

func postgresMakeReading(rows *sql.Rows) thermometers.Reading {
	var name, unit string
	var temp float64
	var readingTime time.Time
	err := rows.Scan(&name, &temp, &unit, &readingTime)
	utils.CheckLog(err)
	return thermometers.Reading{
		Temp: temp,
		Unit: unit,
		Name: name,
		Time: readingTime,
	}
}


func (psql PostgresDB) Latest(thermometer string) (thermometers.Reading, error) {
	return FetchLatest(psql.database, postgresLatest, thermometer, postgresMakeReading)
}

const postgresReadingsBetween = `
	SELECT name, unit, temperature, date  
	FROM readings 
	WHERE name = $1
	AND date >= $2
	AND date <= $3
	ORDER BY id
`
// ALL of the DBS will share a similar between, but they'll differ in how they
// process dates because they'll represent dates differently
func (psql PostgresDB) Between(name string, timestampRange UnixTimestampRange) ([]thermometers.Reading, error) {
	db := psql.database
	tx, stmt := PrepareStatement(db, postgresReadingsBetween)
	defer stmt.Close()
	end := timestampRange.End
	if end == 0 {
		end = int64(time.Now().Unix())
	}
	rows, err := stmt.Query(name, time.Unix(timestampRange.Begin, 0), time.Unix(end, 0))
	defer rows.Close()
	utils.CheckLog(err)
	readings := make([]thermometers.Reading, 0)
	searchError := errors.New("404")
	for rows.Next() {
		searchError = nil
		readings = append(readings, postgresMakeReading(rows))
	}
	tx.Commit()
	return readings, searchError
}
