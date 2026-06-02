package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== 1. Бесконечный цикл + break ===")
	counter := 0
	for {
		fmt.Println("Counter:", counter)
		counter++
		if counter > 3 {
			break
		}
	}

	fmt.Println("\n=== 2. Цикл с условием (эмуляция while) ===")
	counter = 0
	for counter <= 3 {
		fmt.Println("Counter:", counter)
		counter++
	}

	fmt.Println("\n=== 3. Полный цикл for ===")
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	fmt.Println("\n=== 4.1 range по массиву/слайсу ===")
	fruits := []string{"apple", "banana", "cherry"}
	for index, value := range fruits {
		fmt.Printf("Индекс: %d, Значение: %s\n", index, value)
	}

	fmt.Println("\n=== 4.2 range по строке (руны) ===")
	str := "привет"
	for i, r := range str {
		fmt.Printf("Байтовый индекс: %d, Руна: %c\n", i, r)
	}

	fmt.Println("\n=== 4.3 range по map (порядок случаен) ===")
	m := map[string]int{
		"Go":   1,
		"Java": 2,
		"Py":   3,
	}
	for key, value := range m {
		fmt.Printf("Ключ: %s, Значение: %d\n", key, value)
	}

	fmt.Println("\n=== 4.4 range по каналу ===")
	ch := make(chan int)
	// Отправляем данные в канал (в отдельной горутине)
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch) // ВАЖНО: закрываем канал!
	}()
	// Читаем из канала через range
	for value := range ch { // ТОЛЬКО ОДНА переменная!
		fmt.Println("Получено:", value)
	}

	fmt.Println("\n=== 5. Демонстрация проблемы с замыканием ===")
	fmt.Println("Проблема: все горутины видят последнее значение:")
	for i := 1; i <= 3; i++ {
		go func() {
			fmt.Print(i, " ") // Проблема: захват переменной i
		}()
	}
	fmt.Println("\n(Вывод может быть: 4 4 4 или перемешан, но все 4)")

	fmt.Scanln()
}
