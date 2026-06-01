package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	router := mux.NewRouter() // Создание экземпляра роутера

	router.HandleFunc("/C", func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["key"] // Получаем параметр по имени
		fmt.Fprintln(w, id)

	}).Queries("key", "{key:[0-9]+}")
	
	fmt.Println("http://localhost:3000")

	http.ListenAndServe(":3000", router) // Используем реализацию сервера из net/http. gorilla не имеет встроенной реализации

}