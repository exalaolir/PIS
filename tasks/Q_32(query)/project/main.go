package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/Q", handlerQ2).Methods(http.MethodGet).
		Queries("x", "{x}", "y", "{y}")
	r.HandleFunc("/Q", handlerQ3).Methods(http.MethodGet).
		Queries("x", "{x:[0-9]+}")
	r.HandleFunc("/Q", handlerQ1).Methods(http.MethodGet).
		Queries("x", "{x}")

	r.HandleFunc("/Q4", handlerQ4).Methods(http.MethodGet)

	log.Println("Сервер запущен на порту 3000")
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		log.Fatal(err)
	}
}

func handlerQ1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	x := params["x"]
	response := fmt.Sprintf(`{
		"handler":"Q1",
		"x": "%s"
	}`, x)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handlerQ2(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	x := params["x"]
	y := params["y"]

	var xInt, yInt int
	var errX, errY error

	if xInt, errX = strconv.Atoi(x); errX != nil {
		http.Error(w, "ошибка в параметре x: "+errX.Error(), http.StatusBadRequest)
		return
	}
	if yInt, errY = strconv.Atoi(y); errY != nil {
		http.Error(w, "ошибка в параметре y: "+errY.Error(), http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf(`{
		"handler":"Q2",
		"x": "%s",
		"y": "%s",
		"sum": %d
	}`, x, y, xInt+yInt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handlerQ3(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	x := params["x"]

	xInt, _ := strconv.Atoi(x)

	response := fmt.Sprintf(`{
		"handler":"Q3",
		"x": %d,
		"x*x": %d
	}`, xInt, xInt*xInt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handlerQ4(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	x := query.Get("x")
	y := query.Get("y")

	tags := query["tag"]

	response := fmt.Sprintf(`{
		"handler":"Q4",
		"x": "%s",
		"y": "%s",
		"tags": %v
	}`, x, y, tags)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
