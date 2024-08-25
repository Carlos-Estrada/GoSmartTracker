package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var Users []User

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	Users = append(Users, User{ID: 1, Username: "John Doe", Password: hashPassword("password")})
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Password = hashPassword(user.Password)
	Users = append(Users, user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, user := range Users {
		if user.Username == credentials.Username && checkPasswordHash(credentials.Password, user.Password) {
			json.NewEncoder(w).Encode("Logged In Successfully")
			return
		}
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func main() {
	router := mux.NewRouter()
	
	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")

	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}