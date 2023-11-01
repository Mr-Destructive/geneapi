package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/mr-destructive/geneapi/api"
	//"github.com/mr-destructive/geneapi/geneapi"
)

func handleRequests(port int) {
	http.HandleFunc("/", handler.Handler)
	//http.HandleFunc("/register/", geneapi.Register)
	//http.HandleFunc("/generate/", geneapi.Generate)
	//http.HandleFunc("/update/keys/", geneapi.UpdateLLMAPIKeys)
	//log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	port := 8080
	log.Println(port)
	//err := geneapi.InitDB()
	//if err != nil {
	//	log.Fatal(err)
	//}
	handleRequests(port)
}
