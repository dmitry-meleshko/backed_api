package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthPassword obviously should be stored outside the production code
var AuthPassword = "super_secret"

type account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	// sensitive data -- pull from the environment
	envSecret := os.Getenv("AUTH_SECRET")
	if envSecret != "" {
		AuthPassword = envSecret
	}
}

// authHandlerNewToken is an authentication endpoint which will generate a token.
func authHandlerNewToken(w http.ResponseWriter, r *http.Request) {
	var user account

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid username and password combination", http.StatusUnauthorized)
		return
	}

	// TODO: implement authentication data store and validate credentials

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(AuthPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{"token": "` + tokenStr + `"}`)

	return
}

// authValidateHandler is a middleware which ensures that requests are authenticated
var authValidateHandler = func(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// Header format: "Authorization: Bearer TOKEN"
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Missing or invalid auth token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte(AuthPassword), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}
	})
}
