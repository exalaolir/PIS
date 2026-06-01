package main

import (
	"log"
	"net/http"

	"go07_01/rpc"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/rpc", rpc.RpcHandler).Methods("POST")

	log.Println("Сервер прослушивает :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
