package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PostData struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func sayHello(w http.ResponseWriter) {
	fmt.Fprintf(w, "hello\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sayHello(w)
	case http.MethodPost:
		HandlePostRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	var data PostData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Data received successfully", "id": data.ID}
	json.NewEncoder(w).Encode(response)
}
