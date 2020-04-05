#Calor
Calor is a go daemon that takes period readings for a pre-configured list of thermometers, and makes those readings
available via a simple HTTP rest interface.

It's primary purpose was as a method of exploring the go language, and as a way to log the temperatures in my house.
In the hopes that you find it useful for your own projects, it's been written with extensibility in mind, and
adding new thermometer types can be done via implementing the `thermometers.Thermometer` interface.

Everything is in `/internal`, which means Calor does not currently export any packages you can use.
When the code quality comes up a bit, and things are better documented that will change and most things
will move to `/pkg`. Until that time, to use Calor in your own project, you'll have to fork it.
