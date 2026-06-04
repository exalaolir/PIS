package main

//GO: операции (присваивание, бинарные, унарные, логические).
import "fmt"

func main() {
	// 1. Операции присваивания.
	a := 10
	b := 3
	fmt.Println("Обычное присваивание a:", a) // результат: 10
	fmt.Println("Обычное присваивание b:", b) // результат: 3

	a += 5
	fmt.Println("a += 5:", a) // результат: 15

	a -= 2
	fmt.Println("a -= 2:", a) // результат: 13

	a *= 2
	fmt.Println("a *= 2:", a) // результат: 26

	a /= 4
	fmt.Println("a /= 4:", a) // результат: 6

	a %= 3
	fmt.Println("a %= 3:", a) // результат: 0

	// 2. Бинарные арифметические операции.
	x := 10
	y := 4
	fmt.Println("x + y:", x+y) // результат: 14
	fmt.Println("x - y:", x-y) // результат: 6
	fmt.Println("x * y:", x*y) // результат: 40
	fmt.Println("x / y:", x/y) // результат: 2
	fmt.Println("x % y:", x%y) // результат: 2

	// 3. Бинарные операции сравнения.
	fmt.Println("x == y:", x == y) // результат: false
	fmt.Println("x != y:", x != y) // результат: true
	fmt.Println("x > y:", x > y)   // результат: true
	fmt.Println("x < y:", x < y)   // результат: false
	fmt.Println("x >= y:", x >= y) // результат: true
	fmt.Println("x <= y:", x <= y) // результат: false

	// 4. Унарные операции.
	number := 5
	negative := -number
	positive := +number
	fmt.Println("-number:", negative) // результат: -5
	fmt.Println("+number:", positive) // результат: 5

	number++
	fmt.Println("number++:", number) // результат: 6

	number--
	fmt.Println("number--:", number) // результат: 5

	isReady := false
	fmt.Println("!isReady:", !isReady) // результат: true

	// 5. Логические операции.
	age := 20
	hasTicket := true
	isStudent := false

	fmt.Println("age >= 18 && hasTicket:", age >= 18 && hasTicket) // результат: true
	fmt.Println("hasTicket || isStudent:", hasTicket || isStudent) // результат: true
	fmt.Println("!isStudent:", !isStudent)                         // результат: true

	// 6. Бинарные побитовые операции.
	m := 6                       // 6 в двоичной системе: 110
	n := 3                       // 3 в двоичной системе: 011
	fmt.Printf("m & n: %04b\n", m&n)   // результат: 2
	fmt.Println("m | n:", m|n)   // результат: 7
	fmt.Println("m ^ n:", m^n)   // результат: 5
	fmt.Println("m << 1:", m<<1) // результат: 12
	fmt.Println("m >> 1:", m>>1) // результат: 3
}
