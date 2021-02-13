package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 4000

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/gist", showGist)
	mux.HandleFunc("/gist/create", createGist)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
