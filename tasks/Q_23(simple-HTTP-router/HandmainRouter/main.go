package main

import (
	"log"
	"net/http"
)

// Идиоматичный роутер в виде функции с инвертированной логикой (Method -> Path)
func ApiRouter(w http.ResponseWriter, r *http.Request) {
	// Логируем входящий запрос
	log.Printf("[REQ] Входящий запрос: %s %s от %s", r.Method, r.URL.Path, r.RemoteAddr)

	// Устанавливаем общий заголовок для всех ответов (например, JSON)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// ШАГ 1: Маршрутизация по HTTP-методу
	switch r.Method {

	case http.MethodGet:
		// ШАГ 2: Внутри метода распределяем по путям (Только для чтения)
		switch r.URL.Path {
		case "/users":
			w.Write([]byte(`{"action": "get_users", "data": ["Антон", "Иван"]}`))
		case "/status":
			w.Write([]byte(`{"status": "OK", "uptime": "100%"}`))
		default:
			// Если для GET-метода такой путь не найден
			log.Printf("[WARN] GET-путь не найден: %s", r.URL.Path)
			http.Error(w, `{"error": "Ресурс не найден"}`, http.StatusNotFound)
		}

	case http.MethodPost:
		// ШАГ 2: Внутри метода распределяем по путям (Только для создания)
		switch r.URL.Path {
		case "/users":
			// Имитация создания пользователя
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"action": "create_user", "result": "success"}`))
		default:
			log.Printf("[WARN] POST-путь не найден: %s", r.URL.Path)
			http.Error(w, `{"error": "Ресурс не найден"}`, http.StatusNotFound)
		}

	case http.MethodDelete:
		// ШАГ 2: Внутри метода распределяем по путям (Только для удаления)
		switch r.URL.Path {
		case "/users":
			w.Write([]byte(`{"action": "delete_user", "result": "done"}`))
		default:
			http.Error(w, `{"error": "Ресурс не найден"}`, http.StatusNotFound)
		}

	default:
		// Если клиент прислал метод, который мы вообще не обрабатываем (например, PUT или OPTIONS)
		log.Printf("[WARN] Неподдерживаемый метод: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Этот HTTP-метод не поддерживается сервером"}`))
	}
}

func main() {
	// Настраиваем красивый вывод логов (Время и файл)
	log.SetFlags(log.Ltime | log.Lshortfile)

	// Направляем абсолютно все запросы на наш ApiRouter
	http.HandleFunc("/", ApiRouter)

	log.Println("[START] API Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[FATAL] Ошибка запуска: %v", err)
	}
}
