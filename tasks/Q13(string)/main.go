package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Объявление строки.
	name := "Victoria"
	var empty string
	fmt.Println("Имя:", name)
	fmt.Println("Пустая строка:", empty == "")

	// Строки неизменяемые: старую строку не меняем, а создаем новую.
	text := "Hello"
	newText := "h" + text[1:]
	fmt.Println("Новая строка:", newText)

	// len считает байты, а не символы.
	fmt.Println("len(\"Hello\"):", len("Hello"))
	fmt.Println("len(\"Привет\"):", len("Привет"))

	// Индексация строки возвращает байт.
	fmt.Println("Первый байт Hello:", text[0])
	fmt.Printf("Первый символ Hello: %c\n", text[0])

	// Для русских символов удобно использовать []rune.
	russian := "Привет"
	runes := []rune(russian)
	fmt.Println("Количество символов:", len(runes))
	fmt.Println("Первый символ:", string(runes[0]))

	// range проходит по символам rune, но индекс показывает позицию байта.
	for index, char := range russian {
		fmt.Println("range:", index, string(char))
	}

	// Срез английской строки работает просто.
	fmt.Println("Срез Hello[1:4]:", text[1:4])

	// Срез русской строки лучше делать через []rune.
	fmt.Println("Первые 3 символа:", string(runes[0:3]))

	// Склеивание строк через +.
	fullName := name + " Borisova"
	fmt.Println("Склеивание:", fullName)

	// Сравнение строк.
	fmt.Println("Go == Go:", "Go" == "Go")
	fmt.Println("Go == Java:", "Go" == "Java")

	// Популярные функции пакета strings.
	message := "  I learn Go  "
	fmt.Println("TrimSpace:", strings.TrimSpace(message))
	fmt.Println("ToUpper:", strings.ToUpper(message))
	fmt.Println("Contains Go:", strings.Contains(message, "Go"))
	fmt.Println("ReplaceAll:", strings.ReplaceAll(message, "Go", "Golang"))

	// Split разбивает строку, Join склеивает срез строк.
	parts := strings.Split("go,java,python", ",")
	fmt.Println("Split:", parts)
	fmt.Println("Join:", strings.Join(parts, " | "))

	// Преобразование строки в число.
	number, err := strconv.Atoi("123")
	if err == nil {
		fmt.Println("Atoi:", number+10)
	}

	// Преобразование числа в строку.
	numberText := strconv.Itoa(456)
	fmt.Println("Itoa:", "Number: "+numberText)

	// Форматирование строки через fmt.Sprintf.
	info := fmt.Sprintf("Name: %s, Age: %d", name, 25)
	fmt.Println("Sprintf:", info)

	// Преобразование string в []byte и обратно.
	bytes := []byte("Go")
	fmt.Println("[]byte:", bytes)
	fmt.Println("string из []byte:", string(bytes))

	// strings.Builder удобен для сборки строки из частей.
	var builder strings.Builder
	builder.WriteString("Hello")
	builder.WriteString(" ")
	builder.WriteString("Go")
	fmt.Println("Builder:", builder.String())
}
