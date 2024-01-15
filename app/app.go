package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mariosoaresreis/go-hotel/domain/adapters/services"
	"github.com/mariosoaresreis/go-hotel/logger"
	"log"
	"net/http"
	"os"
)

func checkEnvironmentVariables() {
	envProps := []string{
		"ADDRESS",
		"PORT",
		"client_secret",
		"client_id",
		"grant_type",
		"scope",
		"AUTH_URL",
		"OPTII_URL",
	}

	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s not defined. Exiting application...", k))
		}
	}
}
func Start() {
	checkEnvironmentVariables()
	router := mux.NewRouter()
	handlers := JobHandlers{services.NewJobService()}
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/addJob", handlers.addJob)
	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}
