package main

import (
	"GoAssignment/internal/database"
	"GoAssignment/internal/initilizers"

	log "github.com/sirupsen/logrus"
)

func init() {
	initilizers.LoadEnvVariables()
}

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting up Our App")
	_, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	return nil
}
func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
