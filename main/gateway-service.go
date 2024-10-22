package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func fetchUserData(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getUsersGt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userResponse, err := fetchUserData("http://localhost:8001/users")
	if err != nil {
		http.Error(w, "Error fetching user data", http.StatusInternalServerError)
		return
	}
	w.Write(userResponse)
}

func getUserOrderGt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]
	ordersResponse, err := fetchUserData("http://localhost:8002/orders/user/" + userID)
	if err != nil {
		http.Error(w, "Error fetching order data", http.StatusInternalServerError)
		return
	}
	w.Write(ordersResponse)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsersGt).Methods("GET")
	router.HandleFunc("/users/{id}/orders", getUserOrderGt).Methods("GET") // Ensure this calls getUserOrder

	http.ListenAndServe(":8000", router)
}
