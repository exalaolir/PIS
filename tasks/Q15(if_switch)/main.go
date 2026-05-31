package main

import "fmt"

func main() {
	// 1. if: код выполняется, если условие true.
	age := 17
	if age >= 18 {
		fmt.Println("if: доступ разрешен")
	}

	// 2. if else: один блок для true, другой блок для false.
	if age >= 18 {
		fmt.Println("if else: доступ разрешен")
	} else {
		fmt.Println("if else: доступ запрещен") // результат: доступ запрещен
	}

	// 3. else if: несколько условий проверяются сверху вниз.
	score := 75
	if score >= 90 {
		fmt.Println("else if: отлично")
	} else if score >= 60 {
		fmt.Println("else if: зачет") // результат: зачет
	} else {
		fmt.Println("else if: незачет")
	}

	// 4. if с коротким объявлением.
	if number := 10; number%2 == 0 {
		fmt.Println("if с объявлением: число четное") // результат: число четное
	} else {
		fmt.Println("if с объявлением: число нечетное")
	}

	// 5. switch: выбор одного варианта из нескольких.
	day := "Monday"
	switch day {
	case "Monday":
		fmt.Println("switch: понедельник") // результат: понедельник
	case "Tuesday":
		fmt.Println("switch: вторник")
	default:
		fmt.Println("switch: другой день")
	}

	// 6. Несколько значений в одном case.
	role := "admin"
	switch role {
	case "admin", "moderator":
		fmt.Println("case с несколькими значениями: есть права") // результат: есть права
	case "user":
		fmt.Println("case с несколькими значениями: обычный пользователь")
	default:
		fmt.Println("case с несколькими значениями: неизвестная роль")
	}

	// 7. switch без выражения: похож на if else if.
	temperature := 25
	switch {
	case temperature < 0:
		fmt.Println("switch без выражения: мороз")
	case temperature >= 0 && temperature <= 25:
		fmt.Println("switch без выражения: тепло") // результат: тепло
	default:
		fmt.Println("switch без выражения: жарко")
	}

	// 8. switch с коротким объявлением.
	switch month := "May"; month {
	case "December", "January", "February":
		fmt.Println("switch с объявлением: зима")
	case "March", "April", "May":
		fmt.Println("switch с объявлением: весна") // результат: весна
	default:
		fmt.Println("switch с объявлением: другой сезон")
	}

	// 9. fallthrough: специально выполняет следующий case.
	level := 1
	switch level {
	case 1:
		fmt.Println("fallthrough: первый уровень") // результат: первый уровень
		fallthrough
	case 2:
		fmt.Println("fallthrough: второй уровень тоже выполнен") // результат: второй уровень тоже выполнен
	default:
		fmt.Println("fallthrough: другой уровень")
	}

	// 10. type switch: проверка реального типа значения.
	var value any = "Go"
	switch value.(type) {
	case int:
		fmt.Println("type switch: int")
	case string:
		fmt.Println("type switch: string") // результат: string
	case bool:
		fmt.Println("type switch: bool")
	default:
		fmt.Println("type switch: неизвестный тип")
	}
}
