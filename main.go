package main

import (
	"github.com/mariosoaresreis/go-hotel/app"
	"github.com/mariosoaresreis/go-hotel/logger"
)

func main() {
	logger.Info("Initializing application")
	app.Start()
}
