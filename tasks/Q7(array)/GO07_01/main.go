package main

import (
	arrays "GO07_01/Arrays"
)

func main() {
	println("Базовые операции с массивами")
	arrays.DemoBasics()

	println("\nМногомерные массивы")
	arrays.DemoMulti()

	println("\nУказатели и передача массивов")
	arrays.DemoPointers()
}
