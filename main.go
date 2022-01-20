package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/arganaphangquestian/gorest/models"
	"github.com/arganaphangquestian/gorest/utils"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

var USERS = []models.User{}

func init() {
	password, _ := utils.CreateHash("password")
	USERS = append(USERS, []models.User{
		{ID: ksuid.New().String(), Email: "user@mail.com", Password: *password, Name: "User"},
		{ID: ksuid.New().String(), Email: "admin@mail.com", Password: *password, Name: "Admin"},
		{ID: ksuid.New().String(), Email: "another@mail.com", Password: *password, Name: "Another"},
	}...)
}

func main() {
	mux := mux.NewRouter()

	// Index
	mux.HandleFunc("/", index).Methods("GET")
	// Login
	mux.HandleFunc("/login", login).Methods("POST")
	// Protected Router
	mux.HandleFunc("/dashboard", dashboard).Methods("GET")

	fmt.Println("Server is running at PORT 8000")
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")), mux); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&body)
	var user *models.User
	for _, v := range USERS {
		if v.Email == body.Email {
			user = &v
			break
		}
	}
	if user == nil || !utils.CompareHash(user.Password, body.Password) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Email or Password incorrect",
		})
	}
	access_token, _ := utils.CreateToken(user.ID) // Always Success
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Login success",
		"access_token": access_token,
	})
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 || header[0] != "Bearer" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Token doesn't exist",
		})
	}
	userId, err := utils.Verify(header[1])
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Token invalid",
		})
	}
	var user *models.User
	for _, v := range USERS {
		if v.ID == *userId {
			user = &v
			break
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Dashboard Yeay",
		"user":    user,
	})
}
