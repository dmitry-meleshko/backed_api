package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	db = connectToDB()
	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)
	initRoutes(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initRoutes(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Patients API")
	})

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.HandleFunc("/patients", patientHandlerGetAll).Methods("GET")
	api.HandleFunc("/patients/{id:[0-9]+}", patientHandlerGetOne).Methods("GET")
	api.HandleFunc("/patients", patientHandlerAdd).Methods("POST")

	return
}

func connectToDB() *sql.DB {
	server := os.Getenv("GCP_HOSTNAME")
	if server == "" {
		server = "localhost"
	}
	connURL := fmt.Sprintf("postgres://postgres:postgres@%s/project?sslmode=disable", server)
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB via %s: %v", connURL, err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB via %s: %v", connURL, err)
	}
	log.Println("Connected to DB")
	return db
}
