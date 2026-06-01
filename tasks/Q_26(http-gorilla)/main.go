package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	router := mux.NewRouter() // Создание экземпляра роутера

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

	router.HandleFunc(
		"/form",
		showForm,
	).Methods("GET")

	router.HandleFunc(
		"/register",
		register,
	).Methods("POST")

	fmt.Println("http://localhost:3000")

	http.ListenAndServe(":3000", router) // Используем реализацию сервера из net/http. gorilla не имеет встроенной реализации

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func register(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := r.ParseForm()
	if err != nil {
		http.Error(
			w,
			"Invalid form",
			http.StatusBadRequest,
		)
		return
	}

	name := r.PostForm.Get("name")

	age := r.FormValue("age") // Альтернативный способ. Работает без ParseForm, т.к. вызывает этот метод под капотом


	fmt.Fprintf(
		w,
		"Name: %s; age: %s",
		name,
		age,
	)
}

func showForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Registration</title>
	</head>
	<body>
		<h2>Registration Form</h2>

		<form action="/register" method="POST">
			<label>Name:</label>
			<input type="text" name="name">

			<br><br>

			<label>Age:</label>
			<input type="number" name="age">

			<br><br>

			<button type="submit">
				Register
			</button>
		</form>
	</body>
	</html>
	`

	w.Write([]byte(html))
}
