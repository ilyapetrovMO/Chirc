package main

import (
	"log"

	"github.com/ilyapetrovMO/Chirc/internal/options"
	"github.com/ilyapetrovMO/Chirc/internal/users"
)

func main() {
	app := &application{
		Users: &users.Map{},
	}
	app.Logger = log.Default()

	opts := &options.Options{}
	opts.GetOptions(app.Logger)
	app.Options = opts

	err := app.StartAndListen()
	if err != nil {
		app.Logger.Fatalf(err.Error())
	}
}
