package readings

import (
	"database/sql"
	"follis.net/internal/database"
	"follis.net/internal/thermometers"
)

type SqLiteAcceptor struct {
	DBFile string
}

// assumes the DB tables already exist
func (sla SqLiteAcceptor) Accept(reading thermometers.Reading) {
	db, err := sql.Open("sqlite3", sla.DBFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	database.InsertReading(db, reading)
}



