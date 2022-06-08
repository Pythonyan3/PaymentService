package main

import (
	"log"

	"github.com/Pythonyan3/payment-service/internal/app"
)

func main() {
	var app app.Application = app.Application{}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
