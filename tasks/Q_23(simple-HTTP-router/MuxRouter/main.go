package main

import (
	"log"
	"net/http"
)

// ---- 1. ИЗОЛИРОВАННЫЕ ОБРАБОТЧИКИ ДЛЯ КАЖДОГО ДЕЙСТВИЯ ----

// getUsersHandler отвечает ТОЛЬКО за получение списка пользователей
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [%s] Запрос списка пользователей от %s", r.Method, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"action": "get_users", "data": ["Антон", "Иван"]}`))
}

// createUserHandler отвечает ТОЛЬКО за создание пользователя
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [%s] Создание нового пользователя от %s", r.Method, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"action": "create_user", "status": "success"}`))
}

// statusHandler отвечает за проверку работоспособности системы
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] [%s] Проверка статуса сервера от %s", r.Method, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "OK", "version": "1.0.0"}`))
}

// ---- 2. НАСТРОЙКА МАРШРУТИЗАТОРА В MAIN ----

func main() {
	// Настройка красивого формата логов
	log.SetFlags(log.Ltime | log.Lshortfile)

	// Создаем локальный мультиплексор
	mux := http.NewServeMux()

	// ВАЖНО: Указываем МЕТОД перед ПУТЕМ через пробел.
	// Теперь Go сам поймет, куда направить GET, а куда POST.
	mux.HandleFunc("GET /users", getUsersHandler)
	mux.HandleFunc("POST /users", createUserHandler)
	mux.HandleFunc("GET /status", statusHandler)

	log.Println("[START] Сервер успешно запущен на http://localhost:8080")

	// Запуск сервера с нашим кастомным mux
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("[FATAL] Ошибка при работе сервера: %v", err)
	}
}
