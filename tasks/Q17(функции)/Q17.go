package main

import "fmt"

// 1. Функция с параметрами и возвращаемым значением
func add(a, b int) int {
	return a + b
}

// 2. Несколько возвращаемых значений (ВАЖНО!)
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("деление на ноль")
	}
	return a / b, nil
}

// 3. Именованные возвращаемые значения
func getCoordinates() (x, y int) {
	x = 10
	y = 20
	return
}

// 4. Вариативная функция (...int)
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 5. Функция как аргумент (callback)
func apply(x int, f func(int) int) int {
	return f(x)
}

// 6. Замыкание (ВАЖНО ДЛЯ ГОРУТИН!)
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// 7. Функция возвращает функцию
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// 8. Defer (отложенный вызов)
func deferExample() {
	defer fmt.Println("Это последним")
	fmt.Println("Это первым")
}

// 9. Panic/Recover
func safeDiv(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Паника перехвачена:", r)
		}
	}()
	fmt.Println(a / b)
}

// 10. Рекурсия
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	// 1. Обычная функция
	fmt.Println("5+3 =", add(5, 3))

	// 2. Несколько возвращаемых значений + ошибка
	if res, err := divide(10, 2); err == nil {
		fmt.Println("10/2 =", res)
	}

	// 3. Именованные возвращаемые значения
	x, y := getCoordinates()
	fmt.Printf("Координаты: %d, %d\n", x, y)

	// 4. Вариативная функция
	fmt.Println("Сумма:", sum(1, 2, 3, 4, 5))
	nums := []int{10, 20, 30}
	fmt.Println("Сумма слайса:", sum(nums...))

	// 5. Функция в переменной и callback
	square := func(n int) int { return n * n }
	fmt.Println("Квадрат 5:", square(5))
	fmt.Println("apply:", apply(5, square))

	// 6. Замыкание (ВАЖНО!)
	c := counter()
	fmt.Println(c()) // 1
	fmt.Println(c()) // 2

	// 7. Фабрика функций
	double := makeMultiplier(2)
	fmt.Println("double(5):", double(5))

	// 8. Defer
	deferExample()

	// 9. Panic/Recover
	safeDiv(10, 2)
	safeDiv(10, 0) // Не вызовет паники

	// 10. Рекурсия
	fmt.Println("factorial(5):", factorial(5))

	// 11. Анонимная функция с немедленным вызовом
	result := func(a, b int) int { return a * b }(3, 4)
	fmt.Println("3*4 =", result)
}
