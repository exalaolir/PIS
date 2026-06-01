package main

import "fmt"

type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func (c Counter) Value() int {
	return c.value
}

func zero(x *int) {
	*x = 0
}

func appendOne(s []int) {
	s = append(s, 1)
}

func appendOnePtr(s *[]int) {
	*s = append(*s, 1)
}

func swap(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	fmt.Println("Базовые операции с указателями")
	x := 42
	p := &x
	fmt.Printf("x = %d, &x = %p, p = %p, *p = %d\n", x, &x, p, *p)
	*p = 100
	fmt.Printf("После *p = 100: x = %d\n", x)

	fmt.Println("\nnil")
	var nilPtr *int
	if nilPtr == nil {
		fmt.Println("nilPtr равен nil")
	}
	// fmt.Println(*nilPtr) // panic, закомментировано

	fmt.Println("\nПередача указателя в функцию для изменения оригинала")
	a := 5
	fmt.Println("До zero(&a):", a)
	zero(&a)
	fmt.Println("После zero(&a):", a)

	fmt.Println("\nswap через указатели")
	b, c := 1, 2
	fmt.Printf("До swap: b=%d, c=%d\n", b, c)
	swap(&b, &c)
	fmt.Printf("После swap: b=%d, c=%d\n", b, c)

	fmt.Println("\nМетоды с receiver-указателем")
	cnt := Counter{value: 0}
	fmt.Printf("Начальное значение: %d\n", cnt.Value())
	cnt.Inc()
	fmt.Printf("После Inc(): %d\n", cnt.Value())

	fmt.Println("\nУказатели на структуры")
	var p1 *Counter = &Counter{value: 10}
	p2 := new(Counter)
	p3 := &Counter{value: 20}
	fmt.Printf("p1: %+v, p2: %+v, p3: %+v\n", *p1, *p2, *p3)

	fmt.Println("\nslice, map, chan — не требуют указателя для изменения элементов")
	nums := []int{1, 2, 3}
	fmt.Println("Исходный срез:", nums)

	func(s []int) { s[0] = 99 }(nums)
	fmt.Println("После изменения элемента внутри анонимной функции:", nums)

	appendOne(nums)
	fmt.Println("После appendOne (без указателя):", nums)
	appendOnePtr(&nums)
	fmt.Println("После appendOnePtr (с указателем):", nums)

	fmt.Println("\nСравнение указателей и значений")
	v1, v2 := 10, 10
	ptr1, ptr2 := &v1, &v2
	fmt.Printf("ptr1 == ptr2: %t (адреса разные)\n", ptr1 == ptr2)
	fmt.Printf("*ptr1 == *ptr2: %t (значения равны)\n", *ptr1 == *ptr2)
}
