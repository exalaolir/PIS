Проект представляет собой сервер, принимающий POST-запрос с телом типа `multipart/form-data`, читает обычные поля и загруженный файл, сохраняет файл на диск и отправляет ответ с типом `application/json`.

## Общий сценарий работы

1. Клиент отправляет POST запрос с телом на `http://localhost:3000/MP`
2. Роутер (gorilla/mux) направляет запрос (POST `http://localhost:3000/MP`) в функцию `handlerMP`.
3. В `handlerMP` через `ParseMultipartForm` разбирается тело (`net/http`). Извлекаются текстовые поля (`x`, `y`, `s`, `b`) через `FormValue`. Файл (`avatar`) извлекается через `FormFile` и сохраняется в текущей папке `uploads-<имя_файла>`
4. Также в `handlerMP` формируется ответ в JSON формате с полученными данными: передаются `x`, `y`, `s`, `b` с переданными значениями и `file` с именем (`name`) и размером (`size_bytes`) полученного файла.

## Код

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// создание роутера с использованием gorilla/mux
	r := mux.NewRouter()
	// на путь "/MP" с типом метода POST закрепляется обработчик handlerMP
	r.HandleFunc("/MP", handlerMP).Methods(http.MethodPost)

	log.Println("Сервер прослушивает ", 3000)
	// указанием что сервер будет прослушивать localhost:3000
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		// При ошибке (порт занят или другое) логируем ошибку, процесс завершается.
		log.Fatal(err)
	}
}

func handlerMP(w http.ResponseWriter, r *http.Request) {
	// логируем что пришёл запрос на /MP с методом POST
	log.Printf("POST /MP")
	// проверка что multipart/form-data
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Требуется Content-Type: multipart/form-data", http.StatusBadRequest)
		return
	}
	// проверка что данные не превышают 10 Mb. (10 << 20 равноценно с 10 * 2^20 что соответствует 10 Mb)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// иначе через w (http.ResponseWriter) клиенту отправляется ошибка с пояснениями и статусом StatusBadRequest
		http.Error(w, "Не удается разобрать форму: "+err.Error(), http.StatusBadRequest)
	} else {
		// извлекаем текстовые поля с помощью FormValue и название поле и приводим к нужному типу данных через strconv
		x, er1 := strconv.Atoi(r.FormValue("x"))
		y, er2 := strconv.ParseFloat(r.FormValue("y"), 64)
		s := r.FormValue("s")
		b, er3 := strconv.ParseBool(r.FormValue("b"))

		// клиенту отправляется ошибка с пояснениями и статусом StatusBadRequest если с переданными полями что-то не так
		if er1 != nil {
			http.Error(w, "ошибка в параметре x: "+er1.Error(), http.StatusBadRequest)
			return
		}
		if er2 != nil {
			http.Error(w, "ошибка в параметре y: "+er2.Error(), http.StatusBadRequest)
			return
		}
		if er3 != nil {
			http.Error(w, "ошибка в параметре b: "+er3.Error(), http.StatusBadRequest)
			return
		}
		if s == "" {
			http.Error(w, "ошибка в параметре s : параметр обязателен и не может быть пустым", http.StatusBadRequest)
			return
		}

		// извлекаем файл с помощью FormFile и с помощью названия поля.
		// Возвращает FormFile: поток содержимого (io.ReadCloser), *multipart.FileHeader
		// (имя (Filename), размер в байтах (Size), иные заголовки файла (Content-Type, Content-Disposition))
		file, header, err := r.FormFile("avatar")
		if err != nil {
			// если возникли проблемы при извлечении файла (например не передан или не является файлом) отправляем ошибку
			http.Error(w, "файл avatar не загружен: "+err.Error(), http.StatusBadRequest)
			return
		}
		// создаётся/перезаписывается файл через Create и возвращается указатель на новый файл (dst)
		dst, _ := os.Create("./uploads-" + header.Filename)
		// через Copy копируется из file в dst, пока не достигнет EOF
		io.Copy(dst, file)
		// закрываем dst
		dst.Close()
		// закрываем полученный поток содержимого
		file.Close()

		// устанавливаем в ответ заголовок Content-Type в application/json
		w.Header().Set("Content-Type", "application/json")
		// устанавливаем в ответ статус 200
		w.WriteHeader(http.StatusOK)
		// формируем тело ответа
		response := fmt.Sprintf(`{
			"data": {
				"x": "%d",
				"y": "%f",
				"s": "%s",
				"b": "%t",
				"file": {
					"name": "%s",
					"size_bytes": %d
				}
			}
		}`, x, y, s, b, header.Filename, header.Size)
		// в тело ответа записываем ранее сформированное тело и отправляем клиенту
		w.Write([]byte(response))
	}
}