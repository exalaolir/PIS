package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	// Импортируем драйвер анонимно. Он регистрирует себя в пакете database/sql
	_ "modernc.org/sqlite"
)

// User представляет структуру данных в нашем приложении
type User struct {
	ID        int
	Username  string
	Email     *string // Используем указатель для обработки возможных NULL значений
	CreatedAt time.Time
}

func main() {
	// 1. ИНИЦИАЛИЗАЦИЯ И НАСТРОЙКА ПУЛА СОЕДИНЕНИЙ
	// sql.Open НЕ создает соединение с диском/сетью сразу, он только инициализирует пул.
	db, err := sql.Open("sqlite", "./demo.db")
	if err != nil {
		log.Fatalf("Ошибка инициализации пула: %v", err)
	}
	defer db.Close() // Закрываем пул при выходе из программы

	// Настройка параметров пула соединений (крайне важно для Production)
	db.SetMaxOpenConns(10)                 // Максимум 10 одновременно открытых соединений
	db.SetMaxIdleConns(5)                  // Максимум 5 простаивающих соединений в пуле
	db.SetConnMaxLifetime(time.Minute * 5) // Время жизни соединения

	// Реальная проверка связи с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	// Создаем контекст верхнего уровня с таймаутом для безопасности операций
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. СОЗДАНИЕ ТАБЛИЦЫ (Используем ExecContext, так как выборка строк не требуется)
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.ExecContext(ctx, schema)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	// Очистим таблицу перед демонстрацией
	_, _ = db.ExecContext(ctx, "DELETE FROM users")

	// 3. РАБОТА С ТРАНЗАКЦИЯМИ (sql.Tx)
	// Транзакции гарантируют, что все операции пройдут через ОДНО И ТО ЖЕ соединение
	fmt.Println("\n--- Демонстрация Транзакции ---")
	err = demonstrateTransaction(ctx, db)
	if err != nil {
		log.Printf("Транзакция завершилась ошибкой: %v", err)
	}

	// 4. ПОДГОТОВЛЕННЫЕ ВЫРАЖЕНИЯ (Prepared Statements)
	fmt.Println("\n--- Демонстрация Prepared Statements ---")
	err = demonstratePreparedStatement(ctx, db)
	if err != nil {
		log.Printf("Ошибка при работе с Stmt: %v", err)
	}

	// 5. ВЫБОРКА МНОЖЕСТВА СТРОК (QueryContext) И ОБРАБОТКА NULL
	fmt.Println("\n--- Выборка всех пользователей (QueryContext) ---")
	err = queryMultipleRows(ctx, db)
	if err != nil {
		log.Printf("Ошибка при выборке строк: %v", err)
	}

	// 6. ВЫБОРКА ОДНОЙ СТРОКИ (QueryRowContext)
	fmt.Println("\n--- Выборка одного пользователя (QueryRowContext) ---")
	err = querySingleRow(ctx, db, 1) // Ищем существующего ID 1
	if err != nil {
		log.Printf("Ошибка при выборке строки: %v", err)
	}

	// Демонстрация ситуации, когда строка не найдена
	err = querySingleRow(ctx, db, 999) // ID 999 не существует
	if err != nil {
		log.Printf("Ожидаемая обработка отсутствия строки: %v", err)
	}
}

// demonstrateTransaction показывает, как правильно открывать, фиксировать и откатывать транзакции
func demonstrateTransaction(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil) // Второй аргумент — опции (уровень изоляции). nil — по умолчанию.
	if err != nil {
		return err
	}
	// Идиоматичный паттерн: defer Rollback.
	// Если функция завершится по ошибке до Commit(), изменения автоматически откатятся.
	// Если Commit() вызовется успешно, последующий Rollback в defer просто ничего не сделает.
	defer tx.Rollback()

	// Внутри транзакции ВСЕ запросы делаем СТРОГО через объект 'tx', а не 'db'!
	query := "INSERT INTO users (username, email) VALUES (?, ?)"

	// Вставляем пользователя со значением Email (передаем указатель на строку)
	email1 := "Anton@example.com"
	_, err = tx.ExecContext(ctx, query, "Anton", &email1)
	if err != nil {
		return fmt.Errorf("ошибка вставки Anton: %w", err)
	}

	// Вставляем пользователя, у которого Email равен NULL
	_, err = tx.ExecContext(ctx, query, "Bob", nil)
	if err != nil {
		return fmt.Errorf("ошибка вставки Bob: %w", err)
	}

	// Фиксируем транзакцию. Только в этот момент данные физически сохранятся в БД
	if err := tx.Commit(); err != nil {
		return err
	}
	fmt.Println("Транзакция успешно зафиксирована (Commit). Alice и Bob добавлены.")
	return nil
}

// demonstratePreparedStatement показывает эффективное многократное выполнение одного запроса
func demonstratePreparedStatement(ctx context.Context, db *sql.DB) error {
	query := "INSERT INTO users (username, email) VALUES (?, ?)"

	// Компилируем шаблон запроса на стороне СУБД
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close() // Обязательно закрываем statement для освобождения ресурсов!

	// Массив данных для вставки
	usersToInsert := []struct {
		username string
		email    string
	}{
		{"Serg", "Serg@example.com"},
		{"Diana", "diana@example.com"},
	}

	// Выполняем один и тот же stmt с разными параметрами
	for _, u := range usersToInsert {
		res, err := stmt.ExecContext(ctx, u.username, u.email)
		if err != nil {
			return err
		}

		// Демонстрация работы с sql.Result
		lastID, _ := res.LastInsertId()
		rowsAff, _ := res.RowsAffected()
		fmt.Printf("Запросом Stmt добавлен %s (ID: %d), изменено строк: %d\n", u.username, lastID, rowsAff)
	}

	return nil
}

// queryMultipleRows демонстрирует итерацию по результатам выборки и работу с NULL-значениями
func queryMultipleRows(ctx context.Context, db *sql.DB) error {
	// В SQLite плейсхолдером является знак '?', в Postgres это было бы $1, $2
	query := "SELECT id, username, email, created_at FROM users WHERE id > ?"

	rows, err := db.QueryContext(ctx, query, 0)
	if err != nil {
		return err
	}
	// КРИТИЧЕСКИ ВАЖНО: всегда закрывать rows, иначе соединение утечет из пула
	defer rows.Close()

	// Итерируемся по строкам. Цикл завершится, когда закончатся строки
	for rows.Next() {
		var u User

		// Метод Scan копирует значения из текущей строки в переменные Go.
		// Порядок переменных должен СТРОГО соответствовать порядку полей в SELECT.
		// Так как поле 'email' может быть NULL, мы передаем указатель (*string),
		// sql.Scan автоматически запишет nil в указатель, если в БД лежит NULL.
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
		if err != nil {
			return err
		}

		// Красивый вывод с проверкой на NULL
		emailStr := "NULL"
		if u.Email != nil {
			emailStr = *u.Email
		}
		fmt.Printf("ID: %d | Имя: %-8s | Email: %-20s | Создан: %s\n",
			u.ID, u.Username, emailStr, u.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// КРИТИЧЕСКИ ВАЖНО: после завершения цикла проверяем, не прервался ли он из-за ошибки сети
	if err = rows.Err(); err != nil {
		return fmt.Errorf("ошибка во время итерации по строкам: %w", err)
	}

	return nil
}

// querySingleRow демонстрирует оптимальную выборку строго одной строки
func querySingleRow(ctx context.Context, db *sql.DB, userID int) error {
	query := "SELECT id, username, email, created_at FROM users WHERE id = ?"

	var u User
	// QueryRowContext выполняет запрос и возвращает *sql.Row.
	// Метод Scan() на нем нужно вызывать цепочкой. Вызывать Close() здесь не нужно.
	err := db.QueryRowContext(ctx, query, userID).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)

	if err != nil {
		// Специальная проверка: если запись не найдена, возвращается sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("пользователь с ID %d не найден в системе", userID)
		}
		// Любая другая ошибка (проблема с сетью, синтаксис)
		return err
	}

	emailStr := "NULL"
	if u.Email != nil {
		emailStr = *u.Email
	}
	fmt.Printf("Найден пользователь: %s (Email: %s)\n", u.Username, emailStr)
	return nil
}
