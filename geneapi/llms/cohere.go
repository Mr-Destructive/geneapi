package llms

import (
	"log"

	cohere "github.com/cohere-ai/cohere-go"
)

func GenerateCohereAI(request cohere.GenerateOptions, apiKey string) string {
	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		log.Printf("Cohere.CreateClient error: %v\n", err)
		return ""
	}
	resp, err := client.Generate(request)

	if err != nil {
		log.Printf("Cohere.Generation error: %v\n", err)
		return ""
	}
	return resp.Generations[0].Text
}
