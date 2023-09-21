package llms

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func GenerateOpenAI(request openai.CompletionRequest, apiKey string) string {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateCompletion(
		context.Background(),
		request,
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Text
}
