package main // Главный пакет программы.

import ( // Начинаем блок импортов.
	"database/sql" // Пакет для работы с SQL-базами данных.
	f "fmt"        // Алиас импорта: пакет fmt теперь называется f.
	"os"           // Пакет os нужен для чтения строки подключения из переменной окружения.
	. "strings"    // Точечный импорт: функции strings можно вызывать без strings.

	_ "github.com/lib/pq" // Пустой импорт: регистрирует драйвер PostgreSQL для database/sql.
) // Закрываем блок импортов.

type B = int // B является алиасом int.

type A = B // type A = B: A является алиасом типа B.

type C B // type C B: C является новым пользовательским типом на основе B.

func main() { // Точка входа в программу.
	var b B = 10   // Создаем переменную типа B.
	var a A = b    // Присваиваем B в A без преобразования, потому что A - алиас.
	var c C = C(b) // Преобразуем B в C, потому что C - новый тип.

	f.Println("Алиас импорта fmt:", a)                   // Используем fmt через алиас f.
	f.Println("Новый тип C:", c)                         // Выводим значение нового типа C.
	f.Println("Точечный импорт strings:", ToUpper("go")) // Вызываем ToUpper без strings.
	f.Println("Доступные SQL-драйверы:", sql.Drivers())  // Показываем, что драйвер postgres зарегистрирован.

	connectionString := os.Getenv("DATABASE_URL") // Берем строку подключения из переменной окружения DATABASE_URL.
	if connectionString == "" {                   // Проверяем, задана ли переменная окружения.
		connectionString = "postgres://banquet_user:banquet_password@localhost:5432/banquet_constructor?sslmode=disable" // Используем локальный PostgreSQL без Docker.
	} // Завершаем проверку переменной окружения.

	db, err := sql.Open("postgres", connectionString) // Создаем объект подключения через драйвер postgres.
	if err != nil {                                   // Проверяем ошибку создания объекта подключения.
		f.Println("Ошибка sql.Open:", err) // Выводим ошибку, если драйвер не найден.
		return                             // Завершаем программу при ошибке.
	} // Завершаем проверку ошибки.
	defer db.Close() // Закрываем объект подключения в конце программы.

	f.Println("Подключение к PostgreSQL работает")                 // Сообщаем, что подключение успешно проверено.
} // Завершаем функцию main.


//запуск  go run main.go (тут примеры ихсодя из теории)