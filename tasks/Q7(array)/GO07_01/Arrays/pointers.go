package arrays

import "fmt"

func modifyArrayByValue(arr [3]int) {
	arr[0] = 999
}

func modifyArrayByPtr(arr *[3]int) {
	arr[0] = 999
}

func DemoPointers() {

	a := [3]int{1, 2, 3}
	modifyArrayByValue(a)
	fmt.Println("После modifyArrayByValue (копия):", a)

	modifyArrayByPtr(&a)
	fmt.Println("После modifyArrayByPtr (указатель):", a)

	b := [2]int{1, 2}
	c := [2]int{1, 2}
	d := [2]int{2, 3}
	fmt.Println("b == c:", b == c)
	fmt.Println("b == d:", b == d)
}
