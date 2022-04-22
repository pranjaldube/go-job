package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func viperEnvVariable(key string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {
	dbURL := viperEnvVariable("DATABASE_URL")
	DB := InitDB(dbURL)
	h := DBHandler(DB)
	router := mux.NewRouter()

	router.HandleFunc("/api/", h.HealthCheck).Methods("GET")
	router.HandleFunc("/api/submit/", h.JobSubmit).Methods("POST")
	router.HandleFunc("/api/status", h.JobStatus).Methods("GET")

	fmt.Println("Serving at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
