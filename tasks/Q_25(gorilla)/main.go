package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Status:  "OK",
		Message: "GET request received",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "BadRequest",
			Message: "Invalid JSON",
		})
		return
	}
	defer r.Body.Close()

	log.Printf("POST received: Status=%s, Message=%s", req.Status, req.Message)

	response := Response{
		Status:  "OK",
		Message: "POST received: " + req.Message,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:  "BadRequest",
			Message: "Invalid JSON",
		})
		return
	}
	defer r.Body.Close()

	log.Printf("PUT received: Status=%s, Message=%s", req.Status, req.Message)

	response := Response{
		Status:  "OK",
		Message: "PUT request processed",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Status:  "OK",
		Message: "DELETE request processed",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func handleGetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Status:  "OK",
		Message: "User id: " + id,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api", handleGetAll).Methods("GET")
	router.HandleFunc("/api", handlePost).Methods("POST")
	router.HandleFunc("/api", handlePut).Methods("PUT")
	router.HandleFunc("/api", handleDelete).Methods("DELETE")
	router.HandleFunc("/api/{id}", handleGetByID).Methods("GET")

	log.Println("Server start and listen port 3000")

	if err := http.ListenAndServe("localhost:3000", router); err != nil {
		log.Fatal("Server closed:", err)
	}
}
