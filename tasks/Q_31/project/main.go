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
	r := mux.NewRouter()
	r.HandleFunc("/MP", handlerMP).Methods(http.MethodPost)
	log.Println("Сервер прослушивает ", 3000)
	if err := http.ListenAndServe("localhost:3000", r); err != nil {
		log.Fatal(err)
	}
}
func handlerMP(w http.ResponseWriter, r *http.Request) {
	log.Printf("POST /MP")
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		http.Error(w, "Требуется Content-Type: multipart/form-data", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Не удается разобрать форму: "+err.Error(), http.StatusBadRequest)
	} else {
		x, er1 := strconv.Atoi(r.FormValue("x"))
		y, er2 := strconv.ParseFloat(r.FormValue("y"), 64)
		s := r.FormValue("s")
		b, er3 := strconv.ParseBool(r.FormValue("b"))

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

		file, header, err := r.FormFile("avatar")
		if err != nil {
			http.Error(w, "файл avatar не загружен: "+err.Error(), http.StatusBadRequest)
			return
		}

		dst, _ := os.Create("./uploads-" + header.Filename)
		io.Copy(dst, file)
		dst.Close()
		file.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
		w.Write([]byte(response))
	}
}
