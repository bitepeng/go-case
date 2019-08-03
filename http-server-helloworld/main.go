package main

import (
	"fmt"
	"net/http"
)

/**
http://127.0.0.1:8080
*/
func main() {
	// Server
	http.HandleFunc("/", handleHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

// Hello
func handleHello(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World")); err != nil {
		fmt.Println(err)
	}
}
