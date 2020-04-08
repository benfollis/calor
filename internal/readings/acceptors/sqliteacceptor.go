package acceptors

import (
	"calor/internal/database"
	"calor/internal/thermometers"
)

type SqLiteAcceptor struct {
	MyName string
	DB database.CalorDB
}

func (sla SqLiteAcceptor) Name() string {
	return sla.MyName
}
// assumes the DB tables already exist
func (sla SqLiteAcceptor) Accept(reading thermometers.Reading) {
	myDB := sla.DB
	myDB.InsertReading(reading)
}



