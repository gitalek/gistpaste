package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const port = 4000

type ErrBadRequest struct {
	StatusCode int
	message    string
}

func (e ErrBadRequest) Error() string {
	return e.message
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Gistpaste!"))
}

func validateParamId(p string) (int, *ErrBadRequest) {
	if p == "" {
		return 0, &ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			message:    "id query param is not found"}
	}
	id, err := strconv.Atoi(p)
	if err != nil {
		return 0, &ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			message:    "your id query param is not a number [should be positive integer]",
		}
	}
	if id < 1 {
		return 0, &ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			message:    "your id query param is less than or equal to zero [should be more than or equal to one]",
		}
	}
	return id, nil
}

func showGist(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("id")
	id, err := validateParamId(p)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf(w, "Display a specific gist with ID %d", id)
}

func createGist(w http.ResponseWriter, r *http.Request) {
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
