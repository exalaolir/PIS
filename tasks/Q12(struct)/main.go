package main

import (
	"encoding/json"
	"fmt"
)

// Теги структуры для JSON: `json:"city"` и `json:"street"`.
type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

// Вложенная структура: поле Address имеет тип Address.
type User struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}

type Person struct {
	Name string
}

// Встраивание структуры: Person добавлена без имени поля.
type Employee struct {
	Person
	Position string
}

// Структура для сравнения через ==.
type Point struct {
	X int
	Y int
}

// Передача структуры по значению: функция получает копию User.
func changeNameCopy(user User) {
	user.Name = "Alex"
}

// Передача указателя на структуру: функция получает адрес User.
// Изменение оригинала через *User: меняется исходная структура.
func changeNamePointer(user *User) {
	user.Name = "Alex"
}

// Метод структуры: PrintInfo принадлежит типу User.
func (user User) PrintInfo() {
	fmt.Println("Метод PrintInfo:", user.Name, user.Age)
}

func main() {
	// Создание структуры с именами полей.
	user := User{
		Name:  "Anton",
		Age:   25,
		Email: "anton@example.com",
		Address: Address{
			City:   "Minsk",
			Street: "Lenina",
		},
	}

	// Поля структуры и обращение через точку.
	fmt.Println("Поля структуры:", user.Name, user.Age, user.Email)

	// Вложенная структура: обращение через user.Address.City.
	fmt.Println("Вложенная структура:", user.Address.City)

	// Нулевые значения структуры.
	var emptyUser User
	fmt.Printf("Нулевые значения: %+v\n", emptyUser)

	// Передача структуры по значению.
	changeNameCopy(user)
	fmt.Println("После передачи копии:", user.Name)

	// Передача указателя на структуру.
	// Изменение оригинала через *User.
	changeNamePointer(&user)
	fmt.Println("После передачи указателя:", user.Name)

	// Метод структуры.
	user.PrintInfo()

	// Анонимная структура.
	anon := struct {
		Title string
	}{
		Title: "Temporary",
	}
	fmt.Println("Анонимная структура:", anon.Title)

	// Встраивание структуры.
	employee := Employee{
		Person:   Person{Name: "Ivan"},
		Position: "Manager",
	}
	fmt.Println("Встраивание:", employee.Name, employee.Position)

	// Теги структуры для JSON.
	jsonData, _ := json.Marshal(user)
	fmt.Println("JSON с тегами:", string(jsonData))

	// Сравнение структур через ==.
	p1 := Point{X: 1, Y: 2}
	p2 := Point{X: 1, Y: 2}
	fmt.Println("Сравнение структур:", p1 == p2)

	// nil для указателя на структуру.
	var userPointer *User
	fmt.Println("Указатель nil:", userPointer == nil)
}


//все из теории максимально примеров 