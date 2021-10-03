package main

import (
	"log"

	"github.com/ilyapetrovMO/Chirc/internal/options"
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func main() {
	l := log.Default()
	app := &application{
		Logger: l,
		Users:  users.NewMap(l),
	}

	opts := &options.Options{}
	opts.GetOptions(app.Logger)
	app.Options = opts

	err := app.StartAndListen()
	if err != nil {
		app.Logger.Fatalf(err.Error())
	}
}
