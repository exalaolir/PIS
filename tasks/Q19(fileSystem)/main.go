package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// 1. Создаём каталог
	os.MkdirAll("demo", 0755)
	filePath := filepath.Join("demo", "data.txt")

	// 2. Высокоуровневая запись
	os.WriteFile(filePath, []byte("Hello\nLine2\n"), 0644)

	// 3. Высокоуровневое чтение всего файла
	data, _ := os.ReadFile(filePath)
	fmt.Println("Чтение всего файла: ReadFile:\n", string(data))

	// 4. Низкоуровневое чтение с буфером (явное открытие, чтение, закрытие)
	f, _ := os.Open(filePath)
	buf := make([]byte, 4)
	n, _ := f.Read(buf)
	fmt.Printf("Low-level read: %q\n", buf[:n])
	f.Close()

	// 5. Буферизованное построчное чтение
	f, _ = os.Open(filePath)
	scanner := bufio.NewScanner(f)
	fmt.Println("Buffered scan:")
	for scanner.Scan() {
		fmt.Println(" >", scanner.Text())
	}
	f.Close()

	// 6. Добавление в конец
	f, _ = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	f.WriteString("Appended\n")
	f.Close()

	// 7. Копирование через io.Copy
	src, _ := os.Open(filePath)
	dst, _ := os.Create(filepath.Join("demo", "copy.txt"))
	io.Copy(dst, src)
	src.Close()
	dst.Close()

	// 8. Информация о файле
	info, _ := os.Stat(filePath)
	fmt.Printf("File size: %d, isDir: %v\n", info.Size(), info.IsDir())

	// 9. Чтение каталога (высокоуровневое)
	entries, _ := os.ReadDir("demo")
	fmt.Println("Directory entries:")
	for _, e := range entries {
		fmt.Println(" -", e.Name())
	}

	// 10. Рекурсивный обход
	fmt.Println("WalkDir:")
	filepath.WalkDir("demo", func(p string, d os.DirEntry, err error) error {
		fmt.Println("  ", p)
		return nil
	})

	// 11. Удаление
	os.RemoveAll("demo")
	fmt.Println("Removed demo/")
}
