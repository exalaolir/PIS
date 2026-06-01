package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func createUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var user User

	err := json.NewDecoder(
		r.Body,
	).Decode(&user) // чтение json из тела запроса

	if err != nil {
		http.Error(
			w,
			"Invalid JSON",
			http.StatusBadRequest,
		)
		return
	} // Проверка на наличие ошибок сереализации. Если есть ошибки -- возвращаем код 400

	w.Header().Set(
		"Content-Type",
		"application/json",
	) // Установка заголовков 

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} // Пишем json в тело ответа
}

func main() {

	router := mux.NewRouter() // Создание экземпляра роутера

	router.HandleFunc(
		"/users",
		createUser,
	).Methods("POST")

	router.HandleFunc("/A", handler).
		Methods("GET").
		Headers("Content-Type", "text/plain")

	router.HandleFunc("/B/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"] // Получаем параметр по имени
		fmt.Fprintln(w, id)

	}) // Передача параметра, соответствующего реулярному выражению

	router.HandleFunc("/C", func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["key"] // Получаем параметр по имени
		fmt.Fprintln(w, id)

	}).Queries("key", "{key:[0-9]+}")

	fmt.Println("http://localhost:3000")

	http.ListenAndServe(":3000", router) // Используем реализацию сервера из net/http. gorilla не имеет встроенной реализации

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
