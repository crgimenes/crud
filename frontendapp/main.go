package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "teste")
}

func main() {
	fmt.Println("starting frontend app at localhost:5000")

	err := http.ListenAndServe(":5000", http.HandlerFunc(handler))
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
