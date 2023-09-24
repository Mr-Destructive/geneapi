package handler

import (
	"fmt"
	"net/http"

	"github.com/mr-destructive/geneapi/geneapi"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/register/", geneapi.Register)
	http.HandleFunc("/generate/", geneapi.Generate)
	http.HandleFunc("/update/keys/", geneapi.UpdateLLMAPIKeys)
	fmt.Fprintf(w, "Welcome to the LLM API!")
}
