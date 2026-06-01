package main

import (
	"fmt"

	"github.com/fatih/color"

	"q4_struct_application/appinfo"
)

/*
Структура примера:
   Q4(struct_application)
   ├── go.mod
   ├── go.sum
   ├── main.go
   └── appinfo
       └── info.go

Что здесь есть:
   package main - главный пакет исполняемой программы.
   func main() - точка входа, с нее начинается выполнение.
   fmt - встроенный пакет Go.
   github.com/fatih/color - внешний пакет из интернета.
   q4_struct_application/appinfo - наш отдельный пакет из другого файла.

Наглядно про экспорт:
   В appinfo/info.go есть AppName, User, GetMessage.
   Они начинаются с заглавной буквы, поэтому используются здесь:
      appinfo.AppName
      appinfo.User
      appinfo.GetMessage(...)

   В appinfo/info.go есть secretMessage.
   Оно начинается со строчной буквы, поэтому из main.go его вызвать нельзя:
      appinfo.secretMessage() // будет ошибка

Функции и блоки:
   main() - точка входа.
   printInfo() - обычная функция.
   Тело функции - код внутри {}.
   if { ... } - отдельный блок кода.

*/

// C01 - экспортируемая константа внутри package main.
const C01 = 3.14

// main - точка входа в программу.
func main() {
	// Тело функции main начинается здесь.
	user := appinfo.User{Name: "Student"}

	printInfo(user)

	if C01 > 3 {
		// Это блок кода внутри if.
		color.Green("C01 is greater than 3")
	}
}

// printInfo - неэкспортируемая функция package main, потому что с маленькой буквы.
func printInfo(user appinfo.User) {
	// Тело функции printInfo начинается здесь.
	fmt.Println("Exported const:", appinfo.AppName)
	fmt.Println("Exported type field:", user.Name)
	color.Cyan("Exported function: " + appinfo.GetMessage(user))
	fmt.Println("Const from package main:", C01)
}
