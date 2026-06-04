package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Использование: go run main.go <файл>")
		return
	}

	filename := os.Args[1]

	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	// Получаем только имя файла (без пути)
	// Из "C:\бгту\6сем\пис\lr9\go09_01c\test.txt" получаем "test.txt"
	fileNameOnly := filename
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '\\' || filename[i] == '/' {
			fileNameOnly = filename[i+1:]
			break
		}
	}

	// Используем PUT
	url := "http://localhost:3000/" + fileNameOnly

	// Создаем PUT запрос
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		fmt.Printf("Ошибка создания запроса: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка отправки: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Printf("Файл '%s' успешно отправлен на сервер!\n", fileNameOnly)
	} else {
		fmt.Printf("Ошибка: %s\n", resp.Status)
	}
}
