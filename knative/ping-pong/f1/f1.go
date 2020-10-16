package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("SERVERPORT")
	host := os.Getenv("SERVERHOST")
	addr := os.Getenv("SERVERADDR")

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s", addr, port), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Host = host

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error doing request: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	fmt.Fprint(w, buf.String())
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := "8080"

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
