package main

import (
	"fmt"
	"sync"
)

// ПРОБЛЕМА: main завершается раньше горутин
func problem() {
	go fmt.Println("Привет")
	fmt.Println("Пока")
	// Вывод только: Пока (горутина не успела)
}

// РЕШЕНИЕ: WaitGroup
func solution() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Привет из горутины")
	}()

	wg.Wait() // Ждем горутину
	fmt.Println("Программа завершена")
}

// ПАРАЛЛЕЛЬНЫЙ ПОДСЧЕТ
func parallelSum() {
	var wg sync.WaitGroup
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Считаем первую половину
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := 0
		for i := 0; i < 5; i++ {
			sum += numbers[i]
		}
		fmt.Println("Первая часть:", sum)
	}()

	// Считаем вторую половину
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := 0
		for i := 5; i < 10; i++ {
			sum += numbers[i]
		}
		fmt.Println("Вторая часть:", sum)
	}()

	wg.Wait()
}

func main() {
	fmt.Println("=== ПРОБЛЕМА ===")
	problem()

	fmt.Println("\n=== РЕШЕНИЕ ===")
	solution()

	fmt.Println("\n=== ПАРАЛЛЕЛЬНЫЙ ПОДСЧЕТ ===")
	parallelSum()
}
