package arrays

import "fmt"

func DemoMulti() {

	var matrix [2][3]int
	matrix[0][1] = 5
	matrix[1][2] = 9
	fmt.Println("matrix =", matrix)

	tab := [3][4]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
	}
	fmt.Println("tab =", tab)

	fmt.Println("tab[1][2] =", tab[1][2])

	fmt.Println("Перебор вложенными циклами:")
	for i := 0; i < len(tab); i++ {
		for j := 0; j < len(tab[i]); j++ {
			fmt.Printf("%3d ", tab[i][j])
		}
		fmt.Println()
	}

	threeD := [...][2][3]int{
		{{1, 2, 3}, {4, 5, 6}},
		{{7, 8, 9}, {10, 11, 12}},
	}
	fmt.Println("threeD =", threeD)
}
