package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 4000

func home(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("Hello from Gistpaste!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on :%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
