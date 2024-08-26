package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

var Users = struct {
	sync.Mutex
	M map[string]User
}{M: make(map[string]User)}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	bcryptCost := getBcryptCost()

	password, _ := hashPassword("password", bcryptCost)
	Users.M["JohnDoe"] = User{ID: 1, Username: "John Doe", Password: password}
}

func getBcryptCost() int{
	bcCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil || bcCost < 4 || bcCost > 31 {
		return 14
	} 
	return bcCost
}

func hashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

	bcryptCost := getBcryptCost()

	user.Password, err = hashPassword(user.Password, bcryptCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Users.Lock()
	_, exists := Users.M[user.Username]
	if !exists {
		Users.M[user.Username] = user
		Users.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("User created")
	} else {
		Users.Unlock()
		http.Error(w, "Username already exists", http.StatusBadRequest)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Users.Lock()
	user, exists := Users.M[credentials.Username]
	Users.Unlock()
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