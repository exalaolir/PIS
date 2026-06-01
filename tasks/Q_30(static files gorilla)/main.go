package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"github.com/gorilla/mux"
)

const baseDir = "./static"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/static/{filename}", handlerDFS).Methods(http.MethodGet)

	http.ListenAndServe(":3000", r)
}

// GET /static/{filename}
// Объявление функции handlerDFS. Принимает ResponseWriter (для отправки ответа) 
// и Request (содержит информацию о HTTP запросе)
func handlerDFS(w http.ResponseWriter, r *http.Request) {
	
	// vars := mux.Vars(r) - извлекает параметры маршрута из запроса
	// Например, для URL /DFS/static/file.txt, vars["filename"] будет "file.txt"
	vars := mux.Vars(r)
	
	// filename := vars["filename"] - получаем значение параметра "filename" 
	// из извлеченных переменных маршрута
	filename := vars["filename"]
	
	// filePath := filepath.Join(baseDir, filename) - формирует путь к файлу
	// filepath.Join объединяет пути с правильным разделителем для ОС
	// baseDir объявлена глобально как "./static"
	// Например: "./static" + "file.txt" = "./static/file.txt"
	filePath := filepath.Join(baseDir, filename)

	// info, err := os.Stat(filePath) - получает информацию о файле
	// os.Stat возвращает:
	// - info: структуру с информацией о файле (размер, дата, права доступа)
	// - err: ошибку, если файл не существует или недоступен
	info, err := os.Stat(filePath)
	
	// if err == nil - проверяем, что ошибки нет (файл существует и доступен)
	if err == nil {
		
		// if !info.IsDir() - проверяем, что это не директория
		// info.IsDir() возвращает true, если это папка
		// !info.IsDir() - true, если это файл
		if !info.IsDir() {
			
			// downloadName := path.Base(filename) - извлекает имя файла из пути
			// path.Base возвращает последний элемент пути
			// Например: "docs/report.pdf" -> "report.pdf"
			// Используется для безопасного имени файла при скачивании
			downloadName := path.Base(filename)
			
			// w.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
			// Устанавливает HTTP заголовок Content-Disposition
			// "attachment" заставляет браузер скачать файл, а не открывать его
			// filename=... задает имя сохраняемого файла
			w.Header().Set("Content-Disposition", "attachment; filename="+downloadName)
			
			// w.Header().Set("Content-Type", "application/octet-stream")
			// Устанавливает MIME тип содержимого
			// application/octet-stream означает "бинарные данные"
			// Это заставляет браузер обрабатывать файл как бинарный, а не как текст/HTML
			w.Header().Set("Content-Type", "application/octet-stream")
			
			// http.ServeFile(w, r, filePath) - отправляет файл клиенту
			// Эта функция:
			// 1. Открывает файл по пути filePath
			// 2. Читает его содержимое
			// 3. Копирует в ResponseWriter
			// 4. Обрабатывает Range запросы (частичная загрузка)
			// 5. Устанавливает правильные заголовки Content-Length
			http.ServeFile(w, r, filePath)
			
		} else {
			// Если путь ведет к директории, а не к файлу
			// http.Error отправляет HTTP ошибку с сообщением и статус кодом
			// http.StatusBadRequest = 400
			http.Error(w, "not a file", http.StatusBadRequest)
		}
		
	} else if os.IsNotExist(err) {
		// os.IsNotExist(err) - проверяет, является ли ошибка "файл не найден"
		// Если файл не существует, отправляем ошибку 404
		// http.StatusNotFound = 404
		http.Error(w, "file not found", http.StatusNotFound)
		
	} else {
		// Любая другая ошибка (например, нет прав доступа, повреждение диска)
		// Отправляем ошибку сервера 500 с деталями ошибки
		// http.StatusInternalServerError = 500
		// err.Error() преобразует ошибку в строку для отладки
		http.Error(w, "stat error: "+err.Error(), http.StatusInternalServerError)
	}
}