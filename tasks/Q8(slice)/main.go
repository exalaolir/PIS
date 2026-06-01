package main

import "fmt"

func main() {

	s := []int{1, 2, 3}
	fmt.Println("s:", s, "len:", len(s), "cap:", cap(s))

	s = append(s, 4, 5)
	fmt.Println("после добавления:", s)

	sub := s[1:4]
	sub[0] = 99
	fmt.Println("s после sub изменения:", s)

	fmt.Println(s[:2])
	fmt.Println(s[3:])
	fmt.Println(s[1:3])

	dst := make([]int, 3)
	copy(dst, s)
	fmt.Println("скопирован:", dst)

	s = append(s[:2], s[3:]...)
	fmt.Println("после удаления индекса 2:", s)

	var nilSlice []int
	fmt.Println("nilSlice == nil:", nilSlice == nil)
}
