package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type WebDAVServer struct {
	RootDir string
}

func NewWebDAVServer(rootDir string) *WebDAVServer {
	return &WebDAVServer{RootDir: rootDir}
}

func (s *WebDAVServer) getFullPath(path string) string {
	return filepath.Join(s.RootDir, filepath.Clean("/"+path))
}

func (s *WebDAVServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	switch r.Method {
	case "MKCOL":
		s.handleMkcol(w, r)
	case "PUT":
		s.handlePut(w, r)
	case "GET":
		s.handleGet(w, r)
	case "COPY":
		s.handleCopy(w, r)
	case "MOVE":
		s.handleMove(w, r)
	case "DELETE":
		s.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("Неподдерживаемый метод: %s", r.Method)
	}
}

func (s *WebDAVServer) handleMkcol(w http.ResponseWriter, r *http.Request) {
	fullPath := s.getFullPath(r.URL.Path)

	if _, err := os.Stat(fullPath); err == nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("MKCOL: директория уже существует %s", fullPath)
		return
	}

	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Printf("MKCOL ошибка: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("MKCOL: создана директория %s", fullPath)
}

func (s *WebDAVServer) handlePut(w http.ResponseWriter, r *http.Request) {
	fullPath := s.getFullPath(r.URL.Path)

	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Printf("PUT ошибка создания директории: %v", err)
		return
	}

	file, err := os.Create(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("PUT ошибка создания файла: %v", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("PUT ошибка записи: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("PUT: записан файл %s", fullPath)
}

func (s *WebDAVServer) handleGet(w http.ResponseWriter, r *http.Request) {
	fullPath := s.getFullPath(r.URL.Path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("GET ошибка: %v", err)
		return
	}

	if info.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("GET: %s является директорией", fullPath)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("GET ошибка чтения: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
	log.Printf("GET: прочитан файл %s (%d байт)", fullPath, len(data))
}

func (s *WebDAVServer) handleCopy(w http.ResponseWriter, r *http.Request) {
	srcPath := s.getFullPath(r.URL.Path)

	destFull := r.Header.Get("Destination")
	if destFull == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("COPY: отсутствует заголовок Destination")
		return
	}

	destURL, err := url.Parse(destFull)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("COPY ошибка парсинга Destination: %v", err)
		return
	}
	dstPath := s.getFullPath(destURL.Path)

	log.Printf("COPY: %s -> %s", srcPath, dstPath)

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("COPY ошибка: источник не найден %v", err)
		return
	}

	dir := filepath.Dir(dstPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Printf("COPY ошибка создания директории назначения: %v", err)
		return
	}

	if srcInfo.IsDir() {
		err = s.copyDir(srcPath, dstPath)
	} else {
		err = s.copyFile(srcPath, dstPath)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("COPY ошибка копирования: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("COPY: успешно скопирован %s -> %s", srcPath, dstPath)
}

func (s *WebDAVServer) handleMove(w http.ResponseWriter, r *http.Request) {
	srcPath := s.getFullPath(r.URL.Path)

	destFull := r.Header.Get("Destination")
	if destFull == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("MOVE: отсутствует заголовок Destination")
		return
	}

	destURL, err := url.Parse(destFull)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("MOVE ошибка парсинга Destination: %v", err)
		return
	}

	dstPath := s.getFullPath(destURL.Path)

	log.Printf("MOVE: %s -> %s", srcPath, dstPath)
	if _, err := os.Stat(srcPath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("MOVE ошибка: источник не найден %v", err)
		return
	}

	dir := filepath.Dir(dstPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Printf("MOVE ошибка создания директории назначения: %v", err)
		return
	}
	err = os.Rename(srcPath, dstPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("MOVE ошибка перемещения: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("MOVE: успешно перемещён %s -> %s", srcPath, dstPath)
}

func (s *WebDAVServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	fullPath := s.getFullPath(r.URL.Path)

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("DELETE ошибка: %v", err)
		return
	}

	err := os.RemoveAll(fullPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("DELETE ошибка удаления: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE: удалён %s", fullPath)
}

func (s *WebDAVServer) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (s *WebDAVServer) copyDir(src, dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = s.copyDir(srcPath, dstPath)
		} else {
			err = s.copyFile(srcPath, dstPath)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := NewWebDAVServer("./webdav_root")

	err = os.MkdirAll(server.RootDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	port := ":8080"
	log.Printf("WebDAV сервер запущен")
	log.Printf("Адрес: http://localhost%s", port)
	log.Printf("Корневая директория: %s", server.RootDir)
	log.Printf("Поддерживаемые методы: MKCOL, PUT, GET, COPY, MOVE, DELETE")
	log.Printf("Ожидание подключений...")

	err = http.ListenAndServe(port, server)
	if err != nil {
		log.Fatal(err)
	}
}
