package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

/*
   go help

   go version

   Запустить программу без сборки exe:
   go run .

   go run . -name Ivan -repeat 3

   go build -o q2_cli.exe .

   .\q2_cli.exe -name Anna -repeat 2

8. Проверить зависимости и очистить лишнее:
   go mod tidy

   go mod init q2_cli   - создает go.mod

   go mod tidy   -  go.sum создает

*/

func main() {
	name := flag.String("name", "Student", "name to print")
	repeat := flag.Int("repeat", 1, "how many times to print")

	flag.Parse()

	for i := 0; i < *repeat; i++ {
		color.Cyan(fmt.Sprintf("Hello, %s! This is Go CLI demo.", *name))
	}
}
