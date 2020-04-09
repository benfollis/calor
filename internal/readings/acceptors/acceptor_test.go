package acceptors

import (
	"github.com/benfollis/calor/internal/thermometers"
	"time"
)

func ExampleConsoleAcceptor_Accept() {
	reading := thermometers.Reading{
		Temp: 55,
		Unit: "C",
		Name: "Test",
		Time: time.Unix(0, 0),
	}
	ca := ConsoleAcceptor{MyName: "Foo"}
	ca.Accept(reading)
	// Output:
	// Test at time 1969-12-31 16:00:00 -0800 PST of 55 C
}


