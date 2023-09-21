package geneapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Model struct {
	Name        string
	Description string
}

type Request struct {
	Prompt string `json:"prompt"`
}

type Response struct {
	Response string `json:"response"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// Check if the URL has at least 3 parts (including an empty string)
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	AuthenticateUser(w, r)
	llm_name := parts[2]
	modelsHandler(w, r, llm_name)
}

func modelsHandler(w http.ResponseWriter, r *http.Request, llm_name string) {
	var prompt string
	json.NewDecoder(r.Body).Decode(&prompt)
	log.Println(prompt)
	if prompt == "" {
		json.NewEncoder(w).Encode(&Response{
			Response: "Prompt is required",
		})
		return
	}
	req := &Request{
		Prompt: prompt,
	}
	switch llm_name {
	case "openai":
		response := openaiGenerate(req)
		json.NewEncoder(w).Encode(&Response{
			Response: response.Response,
		})
	case "palm2":
		response := palm2Generate(req)
		json.NewEncoder(w).Encode(&Response{
			Response: response.Response,
		})
	}
}

func openaiGenerate(request *Request) *Response {
	return &Response{
		Response: "text-davinci-003",
	}
}

func palm2Generate(request *Request) *Response {
	return &Response{
		Response: "text-bison-001",
	}
}
