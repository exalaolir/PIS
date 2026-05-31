package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/glebarez/sqlite" // Драйвер SQLite (pure Go, без CGO)
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ==========================================
// 1. ОПРЕДЕЛЕНИЕ МОДЕЛЕЙ И СВЯЗЕЙ (СХЕМА БД)
// ==========================================

// Профиль пользователя (Связь One-to-One с User)
type Profile struct {
	ID     uint   `gorm:"primaryKey"`
	Bio    string `gorm:"type:text"`
	UserID uint   // Внешний ключ для User
}

// Заказ (Связь One-to-Many: у одного User может быть много Orders)
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time // Заполнится автоматически
	UserID    uint      // Внешний ключ для User
}

// Язык программирования (Связь Many-to-Many: Пользователи <-> Языки)
type Language struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

// Основная модель Пользователя
type User struct {
	// gorm.Model автоматически добавляет поля ID, CreatedAt, UpdatedAt, DeletedAt (для Soft Delete)
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Age      int    `gorm:"default:18"`
	IsActive bool   `gorm:"default:true"`

	// Определение связей (Ассоциаций)
	Profile   Profile    `gorm:"constraint:OnDelete:CASCADE;"` // One-to-One
	Orders    []Order    `gorm:"constraint:OnDelete:CASCADE;"` // One-to-Many
	Languages []Language `gorm:"many2many:user_languages;"`    // Many-to-Many (создаст промежуточную таблицу)
}

// ==========================================
// 2. ИСПОЛЬЗОВАНИЕ ХУКОВ (HOOKS)
// ==========================================

// BeforeCreate сработает автоматически ПЕРЕД вставкой любого User в базу данных
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf(" [HOOK] Внимание! Собираемся создать пользователя: %s\n", u.Name)
	if u.Name == "Admin" {
		return errors.New("создание пользователя с именем Admin запрещено через API")
	}
	return nil
}

// Глобальная переменная для хранения пула соединений GORM
var db *gorm.DB

func main() {
	var err error

	// ==========================================
	// 3. ИНИЦИАЛИЗАЦИЯ И МИГРАЦИЯ В GORM
	// ==========================================

	// Открываем подключение и включаем подробный вывод SQL-запросов в консоль
	db, err = gorm.Open(sqlite.Open("gorm_demo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Мы будем видеть весь генерируемый SQL!
	})
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	// Авто-миграция: GORM сам создаст таблицы, индексы и внешние ключи по структурам выше
	err = db.AutoMigrate(&User{}, &Profile{}, &Order{}, &Language{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	fmt.Println(" База данных успешно инициализирована и мигрирована!")

	// Настройка маршрутов (Endpoints) нашего демонстрационного сервера
	http.HandleFunc("/create", handleCreate)  // C - Create (Включая связи)
	http.HandleFunc("/read", handleRead)      // R - Read (Фильтры, Сортировки, Preload)
	http.HandleFunc("/update", handleUpdate)  // U - Update (Проблема нуля и map[string]any)
	http.HandleFunc("/delete", handleDelete)  // D - Delete (Soft и Hard Delete)
	http.HandleFunc("/tx", handleTransaction) // Transactions - Транзакции

	fmt.Println(" Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ==========================================
// 4. ОБРАБОТЧИКИ ЗАПРОСОВ (CRUD ЭНДПОИНТЫ)
// ==========================================

// HandleCreate демонстрирует СREATE (Вставка записи вместе со вложенными связями)
// GET http://localhost:8080/create
func handleCreate(w http.ResponseWriter, r *http.Request) {
	// Создаем сложную структуру пользователя со всеми связями прямо в Go
	newUser := User{
		Name: "Alex Go",
		Age:  28,
		Profile: Profile{
			Bio: "Backend Developer love Go & GORM",
		},
		Orders: []Order{
			{Amount: 1500.50},
			{Amount: 99.90},
		},
		Languages: []Language{
			{Name: "Go"},
			{Name: "SQL"},
		},
	}

	// GORM запишет всё одной командой! Он сам поймет порядок вставки,
	// создаст записи во всех таблицах и проставит Foreign Keys (UserID).
	result := db.Create(&newUser)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// После .Create() в newUser.ID уже лежит валидный ID из базы данных!
	w.Write([]byte(fmt.Sprintf("Пользователь создан успешно! ID: %d, Создано строк в структуре: %d", newUser.ID, result.RowsAffected)))
}

// HandleRead демонстрирует READ (Выборка, фильтрация и жадная загрузка связей)
// GET http://localhost:8080/read
func handleRead(w http.ResponseWriter, r *http.Request) {
	var users []User

	// Построение запроса (Method Chaining):
	// 1. Preload загружает связанные данные из других таблиц (Жадная загрузка / Eager Loading)
	// 2. Where добавляет фильтрацию
	// 3. Order сортирует по ID в обратном порядке
	// 4. Find выполняет запрос и складывает результат в слайс users
	err := db.Preload("Profile").
		Preload("Orders").
		Preload("Languages").
		Where("age >= ?", 18).
		Order("id DESC").
		Find(&users).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем текстовый ответ
	var response string
	for _, u := range users {
		response += fmt.Sprintf("Пользователь: %s (Возраст: %d, Активен: %t)\n", u.Name, u.Age, u.IsActive)
		response += fmt.Sprintf("  - Био: %s\n", u.Profile.Bio)
		response += fmt.Sprintf("  - Кол-во заказов: %d\n", len(u.Orders))
		response += "  - Знает языки: "
		for _, lang := range u.Languages {
			response += lang.Name + " "
		}
		response += "\n\n"
	}

	if response == "" {
		response = "База данных пуста. Сначала вызовите /create"
	}

	w.Write([]byte(response))
}

// HandleUpdate демонстрирует UPDATE (Безопасное обновление "нулевых" значений через map/dictionary)
// GET http://localhost:8080/update?id=1
func handleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Укажите ?id= в запросе", http.StatusBadRequest)
		return
	}

	// Инициализируем «всеядный» словарь (Dictionary)
	// Ключ — строго строка (название колонки в БД), значение — любой тип (any)
	updateData := map[string]any{
		"age":       0,     // Мы ХОТИМ обнулить возраст. Если бы передали структуру User{Age: 0}, GORM бы это проигнорировал!
		"is_active": false, // Меняем флаг на false
		"name":      "Updated Alex",
	}

	// Выполняем обновление целевого пользователя
	err := db.Model(&User{}).Where("id = ?", idStr).Updates(updateData).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Пользователь ID %s успешно обновлен через map[string]any! Проверьте /read (Возраст стал 0, Активен стал false)", idStr)))
}

// HandleDelete демонстрирует DELETE (Мягкое и жесткое удаление)
// GET http://localhost:8080/delete?id=1&hard=true
func handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	hardStr := r.URL.Query().Get("hard")

	if idStr == "" {
		http.Error(w, "Укажите ?id= в запросе", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(idStr)
	var user User
	user.ID = uint(id)

	if hardStr == "true" {
		// ЖЕСТКОЕ УДАЛЕНИЕ (Hard Delete) - запись удаляется физически из таблицы с концами
		err := db.Unscoped().Delete(&user).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Пользователь ID %d удален НАВСЕГДА (ФИЗИЧЕСКИ) из базы данных.", id)))
	} else {
		// МЯГКОЕ УДАЛЕНИЕ (Soft Delete) - GORM просто проставит текущее время в поле deleted_at.
		// Запись останется в БД, но обычный .Find() её больше не увидит.
		err := db.Delete(&user).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("Пользователь ID %d удален МЯГКО (Soft Deleted). В БД проставлен deleted_at.", id)))
	}
}

// HandleTransaction демонстрирует TRANSACTIONS (Безопасное выполнение пачки запросов)
// GET http://localhost:8080/tx
func handleTransaction(w http.ResponseWriter, r *http.Request) {
	// Метод db.Transaction автоматически делает BEGIN и COMMIT.
	// Если замыкание возвращает ошибку (err != nil), GORM автоматически сделает ROLLBACK.
	err := db.Transaction(func(tx *gorm.DB) error {

		// Операция 1: Создаем пользователя
		user := User{Name: "Transaction User", Age: 30}
		if err := tx.Create(&user).Error; err != nil {
			return err // Возврат ошибки приведет к откату всей транзакции
		}

		// Операция 2: Создаем для него заказ
		order := Order{Amount: 5000, UserID: user.ID}
		if err := tx.Create(&order).Error; err != nil {
			return err // Если упало тут — первый юзер тоже сотрется из базы!
		}

		// Имитируем непредвиденную ошибку бизнес-логики для проверки отката:
		// Раскомментируйте строчку ниже, чтобы проверить, как транзакция делает ROLLBACK:
		// return errors.New("что-то пошло не так, отменяем все операции!")

		return nil // Ошибок нет — GORM применит изменения (Commit)
	})

	if err != nil {
		http.Error(w, "Транзакция отменена (Rollback): "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Транзакция выполнена успешно! И Юзер, и его Заказ сохранены вместе."))
}
