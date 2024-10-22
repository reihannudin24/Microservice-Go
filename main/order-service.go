package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Order struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Product string `json:"product"`
	Total   int    `json:"total"`
}

var orders = []Order{
	{ID: 1, UserID: 1, Product: "Laptop", Total: 1200},
	{ID: 2, UserID: 2, Product: "Smartphone", Total: 800},
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getUserOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var userOrders []Order
	for _, order := range orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	if len(userOrders) > 0 {
		json.NewEncoder(w).Encode(userOrders)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No orders found for this user"})
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newOrder Order
	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newOrder.ID = len(orders) + 1
	orders = append(orders, newOrder)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/orders", getOrders).Methods("GET")
	router.HandleFunc("/orders/user/{userId}", getUserOrder).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")

	http.ListenAndServe(":8002", router)
}
