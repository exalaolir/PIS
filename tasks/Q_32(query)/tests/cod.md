## Код

```go
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

	//Маршрут с двумя обязательными query-параметрами "x" и "y"
	//Запрос будет принят только если есть ?x=значение&y=значение
	//Значения параметров будут доступны через mux.Vars(r) под ключами "x" и "y"
	r.HandleFunc("/Q", handlerQ2).Methods(http.MethodGet).
		Queries("x", "{x}", "y", "{y}")

    //Маршрут с обязательным query-параметром "x" и валидацией (только цифры)
	r.HandleFunc("/Q", handlerQ3).Methods(http.MethodGet).
		Queries("x", "{x:[0-9]+}")

	//Маршрут с обязательным query-параметром "x". всё что было непринято handlerQ3 будет идти сюда
	r.HandleFunc("/Q", handlerQ1).Methods(http.MethodGet).
		Queries("x", "{x}")


	// 4. Маршрут с опциональными query-параметрами (классический способ через 
	// r.URL.Query())
	//Этот маршрут не использует .Queries(), поэтому принимает любые запросы на /Q4
	r.HandleFunc("/Q4", handlerQ4).Methods(http.MethodGet)

	log.Println("Сервер запущен на порту 3000")
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		log.Fatal(err)
	}
}
// handlerQ1 обрабатывает запросы с одним обязательным query-параметром "x"
func handlerQ1(w http.ResponseWriter, r *http.Request) {
	// Извлекаем query-параметры из mux.Vars (они туда попадают благодаря .Queries())
	params := mux.Vars(r)
	x := params["x"] // получаем значение параметра x
	// Формируем ответ
	response := fmt.Sprintf(`{
		"handler":"Q1",
		"x": "%s"
	}`, x)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// handlerQ2 обрабатывает запросы с двумя обязательными query-параметрами "x" и "y"
func handlerQ2(w http.ResponseWriter, r *http.Request) {
	// Извлекаем оба query-параметра
	params := mux.Vars(r)
	x := params["x"]
	y := params["y"]

	var xInt, yInt int
	var errX, errY error

	// Пробуем преобразовать в числа,
	if xInt, errX = strconv.Atoi(x); errX != nil {
		//если не число возвращает ошибку
		http.Error(w, "ошибка в параметре x: "+errX.Error(), http.StatusBadRequest)
		return
	}
	if yInt, errY = strconv.Atoi(y); errY != nil {
		http.Error(w, "ошибка в параметре y: "+errY.Error(), http.StatusBadRequest)
		return
	}
	// Формируем JSON-ответ
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


// handlerQ3 обрабатывает запросы с параметром "x", который должен содержать только цифры
func handlerQ3(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	x := params["x"] // здесь x гарантированно состоит только из цифр (благодаря regexp [0-9]+)

	// Преобразуем в число
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

// handlerQ4 обрабатывает запросы с опциональными query-параметрами (классический способ)
func handlerQ4(w http.ResponseWriter, r *http.Request) {
	// Классический способ извлечения query-параметров (через r.URL.Query())
	query := r.URL.Query()

	// Get возвращает значение первого параметра с указанным ключом
	// или пустую строку, если параметр отсутствует
	x := query.Get("x")
	y := query.Get("y")

	// Можно получить все значения, если параметр повторяется несколько раз
	// Например: /Q4?tag=go&tag=web&tag=api
	tags := query["tag"] // вернет []string{"go", "web", "api"}

	// Формируем ответ
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