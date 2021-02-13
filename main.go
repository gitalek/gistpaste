package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = 4000

func home(w http.ResponseWriter, r *http.Request)  {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Gistpaste!"))
}

func showGist (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a specific gist..."))
}

func createGist (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		statusCode := 405
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
    w.Write([]byte("Create a new gist..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/gist", showGist)
	mux.HandleFunc("/gist/create", createGist)

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on :%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
