package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

var users = struct {
	sync.Mutex
	m map[string]User
}{m: make(map[string]User)}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	defaultUser := User{ID: 1, Username: "John Doe"}
	defaultUser.Password, _ = hashPassword("password", getBcryptCost())
	users.m[defaultUser.Username] = defaultUser
}

func getBcryptCost() int {
	bcCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil || bcCost < bcrypt.MinCost || bcCost > bcrypt.MaxCost {
		return bcrypt.DefaultCost
	}
	return bcCost
}

func hashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password, err := hashPassword(user.Password, getBcryptCost())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users.Lock()
	defer users.Unlock()
	if _, exists := users.m[user.Username]; !exists {
		users.m[user.Username] = user
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User created")
	} else {
		http.Error(w, "Username already exists", http.StatusBadRequest)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users.Lock()
	user, exists := users.m[credentials.Username]
	users.Unlock()
	if exists && checkPasswordHash(credentials.Password, user.Password) {
		json.NewEncoder(w).Encode("Logged In Successfully")
		return
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")

	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}