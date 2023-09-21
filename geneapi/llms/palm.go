package llms

import (
	"fmt"

	palm "github.com/mr-destructive/palm"
)

func GeneratePaLM2(request palm.PromptConfig, apiKey string) string {
	client := palm.NewClient(apiKey)
	prompt := request.Prompt.Text
	resp, err := client.ChatPrompt(prompt)
    fmt.Println(resp)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Candidates[0].Output
}
