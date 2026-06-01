package main

import (
	"fmt"
	"os"
)

/*


Зачем нужен config file:
   Настройки можно менять в config.txt без изменения кода программы.
*/

func main() {
	data, err := os.ReadFile("config.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(data))
}
