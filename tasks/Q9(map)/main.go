package main

import "fmt"

func main() {

	scores := make(map[string]int)

	scores["Андрей"] = 100
	scores["Вика"] = 90
	scores["Даша"] = 85

	scores["Андрей"] = 95

	if val, ok := scores["Антон"]; ok {
		fmt.Println("Счет Антона:", val)
	} else {
		fmt.Println("Антон не найден")
	}

	delete(scores, "Вика")

	fmt.Println("Scores:")
	for name, score := range scores {
		fmt.Println(name, score)
	}

	var nilMap map[string]int
	fmt.Println("nilMap is nil?", nilMap == nil)
	nilMap = make(map[string]int)
	nilMap["test"] = 123
	fmt.Println(nilMap)
}
