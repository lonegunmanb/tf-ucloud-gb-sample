package main

import (
	"fmt"
	"log"
	"net/http"
)

const version = 1

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %d", version)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
