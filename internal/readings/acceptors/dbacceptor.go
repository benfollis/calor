package acceptors

import (
	"github.com/benfollis/calor/internal/database"
	"github.com/benfollis/calor/internal/thermometers"
)

// A DBAcceptor emits the readings to a CalorDB
type DBAcceptor struct {
	MyName string
	DB database.CalorDB
}

func (dba DBAcceptor) Name() string {
	return dba.MyName
}
// assumes the DB tables already exist
func (dba DBAcceptor) Accept(reading thermometers.Reading) {
	myDB := dba.DB
	myDB.InsertReading(reading)
}



