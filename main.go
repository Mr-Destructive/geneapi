package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mr-destructive/geneapi/geneapi"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the LLM API!")
}

func handleRequests(port int) {
	http.HandleFunc("/", index)
	http.HandleFunc("/register/", geneapi.Register)
	http.HandleFunc("/generate/", geneapi.Generate)
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	port := 8080
	log.Println(port)
	dbPath := "geneapi/db.sqlite3"
	log.Println(dbPath)
	err := geneapi.InitDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	handleRequests(port)
}
