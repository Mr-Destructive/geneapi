package llms

import (
	"log"

	"github.com/3JoB/anthropic-sdk-go"
)

func GenerateAnthropicAI(request anthropic.Opts, apiKey string) string {
	client, err := anthropic.New(apiKey, "")
	if err != nil {
		log.Printf("anthropic.client error: %v\n", err)
		return ""
	}
	resp, err := client.Send(&request)

	if err != nil {
		log.Printf("Cohere.Generation error: %v\n", err)
		return ""
	}
	return resp.Response.Completion
}
