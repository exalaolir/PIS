package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// UserRequest описывает структуру входящего JSON-запроса.
// Теги `json:"..."` обязательны, чтобы Go понимал, какие поля JSON соответствуют полям структуры.
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// UserResponse описывает структуру исходящего JSON-ответа.
type UserResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// handleUser — функция-обработчик для маршрута /api/user
func handleUser(w http.ResponseWriter, r *http.Request) {
	// 1. Проверка метода. Наш эндпоинт должен принимать ТОЛЬКО POST.
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed) // 405 Method Not Allowed
		_, _ = w.Write([]byte(`{"error": "Разрешен только метод POST"}`))
		return
	}

	// 2. Декодирование JSON из тела запроса (r.Body) в структуру UserRequest
	var req UserRequest

	// Используем json.NewDecoder, так как мы читаем из потока (io.Reader)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		_, _ = w.Write([]byte(`{"error": "Некорректный формат JSON"}`))
		return
	}
	defer r.Body.Close() // Закрываем тело запроса после чтения

	// 3. Простейшая валидация данных
	if req.Name == "" || req.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422 Unprocessable Entity
		_, _ = w.Write([]byte(`{"error": "Поля 'name' и 'email' обязательны для заполнения"}`))
		return
	}

	// Логируем полученные данные на стороне сервера
	fmt.Printf("[LOG] Успешно получен пользователь: %s (%s), возраст: %d\n", req.Name, req.Email, req.Age)

	// 4. Формирование успешного ответа
	response := UserResponse{
		Status:    "success",
		Message:   fmt.Sprintf("Пользователь %s успешно зарегистрирован!", req.Name),
		CreatedAt: time.Now(),
	}

	// Устанавливаем заголовок Content-Type, чтобы клиент знал, что ему возвращают JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	// Кодируем структуру в JSON прямо в ResponseWriter
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Если что-то пошло не так при кодировании ответа
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}

func main() {
	// Регистрируем наш обработчик для пути "/api/user"
	http.HandleFunc("/api/user", handleUser)

	// Определяем порт
	port := ":8080"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)

	// Запускаем сервер. Второй аргумент nil означает использование стандартного маршрутизатора (DefaultServeMux)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
