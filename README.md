#Calor
Calor is a go daemon that takes period readings for a pre-configured list of thermometers, and makes those readings
available via a simple HTTP rest interface.

It's primary purpose was as a method of exploring the go language, and as a way to log the temperatures in my house.
In the hopes that you find it useful for your own projects, it's been written with extensibility in mind, and
adding new thermometer types can be done via implementing the `thermometers.Thermometer` interface.

Everything is in `/internal`, which means Calor does not currently export any packages you can use.
When the code quality comes up a bit, and things are better documented that will change and most things
will move to `/pkg`. Until that time, to use Calor in your own project, you'll have to fork it.


To build and run the commands are
```
go build cmd/calor/calor.go
./calor --config <Full path to your config file>
```
Alternately, if you have a config file in `/etc/calor/calor.json` already, then it's just
```
go build cmd/calor/calor.go
./calor
```

The config file calor is a JSON document with 4 keys
1. `Database` The description of your database configuration, which is used to store the readings from the thermometers.
Currently only SQLite is supported. See the sample config for more details.
1. `ReadAcceptors` A JSON array of Read Acceptors. Read Acceptors are the things that accept the readings from the thermometers
Currently there is a Console read acceptor, that just outputs the readings to the console, and the Sqlite read acceptor,
which writes them to the Sqlite DB configured in the database section
1. `Thermometers` A JSON array of thermometer configurations. See the sample_calor_config.json for some example
types. The `UpdateInterval` option is an integer and is the number of seconds you want a reading taken from that thermometer.
1. `Port` this is the port you want the web server to run on. It does not support HTTPS

# Init scripts
A basic script for systemd is contained in `init`