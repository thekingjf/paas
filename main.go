package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/health", health)
	http.ListenAndServe(":8080", nil)

}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}
