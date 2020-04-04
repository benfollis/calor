package acceptors

import (
	"database/sql"
	"follis.net/internal/database"
	"follis.net/internal/thermometers"
)

type SqLiteAcceptor struct {
	MyName string
	DBFile string
}

func (sla SqLiteAcceptor) Name() string {
	return sla.MyName
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



