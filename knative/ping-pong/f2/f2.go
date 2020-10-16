package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "f2 said hi!\n")
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := "8080"

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
