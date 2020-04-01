package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var _router *mux.Router
var _fakeID int

func TestMain(m *testing.M) {
	db = connectToDB()
	defer db.Close()

	code := m.Run()

	// clean up fake record - ID was set during tests
	p := &patientItem{ID: _fakeID}
	if err := p.deletePatient(db); err != nil {
		panic(err)
	}

	os.Exit(code)

}

func TestAuth(t *testing.T) {
	jsonStr := []byte(`{"username":"test_user","password":"test_pass"}`)
	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("An error occurred. %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/auth", authHandlerNewToken).Methods("POST")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Response code %d doesn't match %d\n", w.Code, http.StatusOK)
		return
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Missing JSON header")
		return
	}

	expected := string(`{\\"token\\": \\".+\\"}`)
	if matched, _ := regexp.MatchString(expected, string(body)); !matched {
		t.Errorf("Response body differs")
		return
	}
}

func TestPatientsAllAuth(t *testing.T) {
	req, err := http.NewRequest("GET", "/patients", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
		return
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter()

	// authenticated endpoint
	router.HandleFunc("/patients", authValidateHandler(patientHandlerGetAll)).Methods("GET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "test_user",
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(AuthPassword))
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code %d doesn't match %d\n", w.Code, http.StatusOK)
		return
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Missing JSON header")
		return
	}

	expected := string(`\[{"id":\d+,"first_name":.+`)
	if matched, _ := regexp.MatchString(expected, string(body)); !matched {
		t.Errorf("Response body differs")
		return
	}
}

func TestPatientMissing(t *testing.T) {
	req, err := http.NewRequest("GET", "/patients/0", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
		return
	}
	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/patients/{id:[0-9]+}", patientHandlerGetOne).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Response code %d doesn't match %d\n", w.Code, http.StatusNotFound)
		return
	}
}

func TestPatientAdd(t *testing.T) {
	jsonStr := []byte(`{"first_name":"ONLY FOR","last_name":"TESTING"}`)
	req, err := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("An error occurred. %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	http.HandlerFunc(patientHandlerAdd).ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if w.Code != http.StatusCreated {
		t.Errorf("Response code %d doesn't match %d\n", w.Code, http.StatusNotFound)
		return
	}
	re := regexp.MustCompile(`"Created patient ID: (\d+)"`)
	match := re.FindStringSubmatch(string(body))
	_fakeID, err = strconv.Atoi(match[1])
	if err != nil {
		t.Errorf("Invalid Patient ID returned: %v", match)
		return
	}

}

func TestPatientPresent(t *testing.T) {
	req, err := http.NewRequest("GET", "/patients/"+strconv.Itoa(_fakeID), nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
		return
	}

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/patients/{id:[0-9]+}", patientHandlerGetOne).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code %d doesn't match %d\n", w.Code, http.StatusOK)
		return
	}

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Missing JSON header")
		return
	}

	expected := string(`\{"id":\d+,"first_name":.+`)
	if matched, _ := regexp.MatchString(expected, string(body)); !matched {
		t.Errorf("Response body differs")
		return
	}
}
