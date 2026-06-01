package arrays

import "fmt"

func DemoBasics() {

	var a [3]int
	a[0] = 10
	a[1] = 20
	a[2] = 30
	fmt.Println("a =", a)

	b := [3]int{1, 2, 3}
	fmt.Println("b =", b)

	c := [5]int{1: 100, 4: 400}
	fmt.Println("c =", c)

	fmt.Print("Автоопределение длины: ")
	d := [...]int{5, 6, 7, 8}
	fmt.Printf("d = %v, длина = %d\n", d, len(d))

	fmt.Print("Копирование массива: ")
	e := d
	e[0] = 999
	fmt.Println("d =", d, "(не изменился)")
	fmt.Println("e =", e)

	fmt.Print("Обход по индексу: ")
	for i := 0; i < len(b); i++ {
		fmt.Printf("%d ", b[i])
	}
	fmt.Println()

	fmt.Print("Обход через range: ")
	for idx, val := range b {
		fmt.Printf("[%d]=%d ", idx, val)
	}
	fmt.Println()
}
