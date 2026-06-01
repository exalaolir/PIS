package main

import "fmt"

/*
Go - язык со строгой типизацией.
Тип можно указать явно, а можно дать Go вывести тип автоматически.

Константы вычисляются на этапе компиляции.
Они не могут зависеть от значений, которые появляются только во время выполнения.
*/

const Pi float32 = 3.14159          // типизированная константа
const Greeting = "Hello from const" // нетипизированная константа
const RussianText = "Привет"        // русские буквы занимают больше байт в UTF-8

func main() {
	var a int = 10
	var a8 int8 = 8
	var b64 float64 = 3.14
	var ok bool = true
	var text string = "Hello"
	var symbol rune = 'A' //символ Unicode, внутри это int32
	var oneByte byte = 255

	autoType := 100 // короткое объявление, тип int выводится автоматически

	fmt.Println("int:", a)
	fmt.Println("int8:", a8)
	fmt.Println("float64:", b64)
	fmt.Println("bool:", ok)
	fmt.Println("string:", text)
	fmt.Println("rune:", symbol)
	fmt.Println("byte:", oneByte)
	fmt.Println("auto type:", autoType)
	fmt.Printf("auto type type: %T\n", autoType)
//%T — специальный шаблон, который выводит тип значения.
	fmt.Println("Pi:", Pi)
	fmt.Printf("Pi type: %T\n", Pi)
	fmt.Println("Greeting:", Greeting)
	fmt.Printf("Greeting type: %T\n", Greeting)

	fmt.Println("len Hello:", len("Hello"))      // 5 байт
	fmt.Println("len Привет:", len(RussianText)) // больше, потому что UTF-8
}
