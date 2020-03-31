package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func patientHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting all patients")
	ps := patients{}
	pItems, err := ps.getPatients(db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			http.Error(w, "No patient records found", http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(pItems)
	w.Write(resp)
}

func patientHandlerGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Patient ID", http.StatusBadRequest)
		return
	}

	p := &patientItem{ID: id}
	log.Println("Getting patient ID: ", p)

	if err := p.getPatient(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			http.Error(w, "Patient not found. ID: "+strconv.Itoa(id), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(p)
	w.Write(resp)

	return
}

func patientHandlerAdd(w http.ResponseWriter, r *http.Request) {
	// extract patient JSON object from POSt request
	var p patientItem
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	log.Println("Adding new patient")
	if err := p.addPatient(db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal("Created patient ID: " + strconv.Itoa(p.ID))
	w.Write(resp)

}

// TODO: implement handlers for Physician and Visit methods
