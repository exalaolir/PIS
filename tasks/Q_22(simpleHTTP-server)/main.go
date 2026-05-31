package main

import (
	"fmt"
	"net/http"
	"time"
)

// helloHandler — функция-обработчик для маршрута "/hello"
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Вывод в консоль сервера (стандартный поток вывода)
	fmt.Printf("[НАЧАЛО] Получен запрос от %s к %s\n", r.RemoteAddr, r.URL.Path)

	// Искусственная задержка в 3 секунды для имитации тяжелой задачи
	// Это поможет нам увидеть, что сервер не блокируется одним запросом
	time.Sleep(3 * time.Second)

	// 2. Отправка ответа клиенту (браузеру/curl)
	fmt.Fprintf(w, "Привет! Твой запрос успешно обработан.\n")

	// 3. Вывод в консоль об успешном завершении
	fmt.Printf("[КОНЕЦ] Запрос от %s обработан\n", r.RemoteAddr)
}

func main() {
	// Регистрируем обработчик helloHandler для пути "/hello"
	// Используется DefaultServeMux (стандартный маршрутизатор)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	fmt.Println("Для проверки откройте несколько вкладок в браузере или используйте curl.")

	// Запускаем сервер на порту 8080.
	// Второй параметр nil означает, что используется DefaultServeMux.
	// Эта функция блокирует выполнение main, пока сервер не будет остановлен (или не произойдет ошибка).
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}
