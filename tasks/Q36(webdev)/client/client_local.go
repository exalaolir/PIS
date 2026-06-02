package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type WebDAVClient struct {
	BaseURL  string
	Username string //
	Password string //
	Client   *http.Client
}

func NewWebDAVClient(baseURL string) *WebDAVClient {
	return &WebDAVClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *WebDAVClient) doRequest(method, urlPath string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.BaseURL + urlPath
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	return c.Client.Do(req)
}

func (c *WebDAVClient) Mkcol(dirPath string) error {
	resp, err := c.doRequest("MKCOL", dirPath, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 || resp.StatusCode == 405 {
		return nil
	}
	return fmt.Errorf("MKCOL failed: %s", resp.Status)
}

func (c *WebDAVClient) Put(filePath string, data []byte) error {
	reader := bytes.NewReader(data)
	resp, err := c.doRequest("PUT", filePath, reader, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("PUT failed: %s", resp.Status)
}

func (c *WebDAVClient) Get(filePath string) ([]byte, error) {
	resp, err := c.doRequest("GET", filePath, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET failed: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func (c *WebDAVClient) Copy(srcPath, dstPath string) error {
	headers := map[string]string{
		"Destination": c.BaseURL + dstPath,
	}
	resp, err := c.doRequest("COPY", srcPath, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("COPY failed: %s", resp.Status)
}

func (c *WebDAVClient) Move(srcPath, dstPath string) error {
	headers := map[string]string{
		"Destination": c.BaseURL + dstPath,
	}
	resp, err := c.doRequest("MOVE", srcPath, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("MOVE failed: %s", resp.Status)
}

func (c *WebDAVClient) Delete(path string) error {
	resp, err := c.doRequest("DELETE", path, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		return nil
	}
	return fmt.Errorf("DELETE failed: %s", resp.Status)
}

func main() {
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := NewWebDAVClient("http://localhost:8080")

	log.Println("WebDAV клиент тестирование локального сервера")

	err = client.Mkcol("/test_folder")
	if err != nil {
		log.Printf("MKCOL ошибка: %v", err)
	} else {
		log.Println("MKCOL: директория /test_folder создана")
	}

	err = client.Put("/test_folder/file.txt", []byte("Hello from WebDAV client!"))
	if err != nil {
		log.Printf("PUT ошибка: %v", err)
	} else {
		log.Println("PUT: файл /test_folder/file.txt записан")
	}

	data, err := client.Get("/test_folder/file.txt")
	if err != nil {
		log.Printf("GET ошибка: %v", err)
	} else {
		log.Printf("GET: содержимое файла: %s", string(data))
	}

	err = client.Copy("/test_folder/file.txt", "/test_folder/copy.txt")
	if err != nil {
		log.Printf("COPY ошибка: %v", err)
	} else {
		log.Println("COPY: файл скопирован в /test_folder/copy.txt")
	}

	err = client.Move("/test_folder/copy.txt", "/test_folder/moved.txt")
	if err != nil {
		log.Printf("MOVE ошибка: %v", err)
	} else {
		log.Println("MOVE: файл перемещен в /test_folder/moved.txt")
	}

	err = client.Delete("/test_folder/moved.txt")
	if err != nil {
		log.Printf("DELETE ошибка: %v", err)
	} else {
		log.Println("DELETE: файл /test_folder/moved.txt удален")
	}

	err = client.Delete("/test_folder")
	if err != nil {
		log.Printf("DELETE директории ошибка: %v", err)
	} else {
		log.Println("DELETE: директория /test_folder удалена")
	}

	log.Println("Тестирование завершено")
}
