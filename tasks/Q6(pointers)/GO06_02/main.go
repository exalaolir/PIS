package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, web programmers!")
}

func main() {
	http.HandleFunc("/", handler)
	//go func() {
	http.ListenAndServe("localhost:3000", nil)
	//}()
	fmt.Println("main")
}
