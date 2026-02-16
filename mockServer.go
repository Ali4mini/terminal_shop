package main

import (
	"fmt"
	"net/http"
)

func StartMockServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "arch;ubuntu;manjaro;debian;fedora;eye-for-eye;nas")
	})

	fmt.Printf("listening on port 9991 \n")
	http.ListenAndServe("localhost:9991", mux)

}
