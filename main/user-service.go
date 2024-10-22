package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users = []User{
		{ID: 1, Name: "Reihan", Email: "reihannudin24@gmail.com"},
		{ID: 2, Name: "Jenni", Email: "jenni24@gmail.com"},
		{ID: 3, Name: "Shannon", Email: "shannon24@gmail.com"},
		{ID: 4, Name: "Nicole", Email: "nicole24@gmail.com"},
	}
	mutex = &sync.Mutex{}
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "User not found"})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")

	http.ListenAndServe(":8001", router)
}
