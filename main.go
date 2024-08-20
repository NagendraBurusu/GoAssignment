package main

import (
	"GoAssignment/internal/database"
	"GoAssignment/internal/initilizers"
	"GoAssignment/internal/student"
	transportHTTP "GoAssignment/internal/transport"

	log "github.com/sirupsen/logrus"
)

func init() {
	initilizers.LoadEnvVariables()
	initilizers.GetLogs()
}

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting up Our App")
	db, err := database.NewDatabase()

	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	stdentService := student.NewService(db)
	handler := transportHTTP.NewHandler(stdentService)
	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}
	return nil
}
func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
