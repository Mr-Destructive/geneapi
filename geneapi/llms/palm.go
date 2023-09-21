package llms

import (
	"log"

	palm "github.com/mr-destructive/palm"
)

func GeneratePaLM2(request palm.PromptConfig, apiKey string) string {
	client := palm.NewClient(apiKey)
	prompt := request.Prompt.Text
	if client == nil {
		log.Printf("Could not create client, invalid API key\n")
		return ""
	}
	resp, err := client.ChatPrompt(prompt)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return ""
	}
	return resp.Candidates[0].Output
}
