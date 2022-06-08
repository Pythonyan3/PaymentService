package app

import (
	"log"
	"os"
	"os/signal"
)

type Application struct{}

func (app *Application) Run() error {
	// application entrypoint method, initialize whole things.
	log.Println("Service is started!")

	// waiting for Ctrl + C to exit application
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Service is shutted down!")

	return nil
}
