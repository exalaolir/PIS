package main

import (
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func main() {
	handler := &webdav.Handler{
		Prefix:     "/",
		FileSystem: webdav.Dir("./data"),
		LockSystem: webdav.NewMemLS(),
	}

	log.Println("WebDAV server on :3000")
	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Fatal(err)
	}

}
