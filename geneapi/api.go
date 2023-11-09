package geneapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/3JoB/anthropic-sdk-go"
	cohere "github.com/cohere-ai/cohere-go"
	"github.com/mr-destructive/geneapi/geneapi/llms"
	"github.com/mr-destructive/geneapi/geneapi/types"
	palm "github.com/mr-destructive/palm"
	openai "github.com/sashabaranov/go-openai"
)

func GeneAI(req types.Request, llmName string, llmKeys map[string]string) (string, error) {
	return geneHandler(req, llmName, llmKeys)
}
func Generate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	// Check if the URL has at least 3 parts (including an empty string)
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	llmKeys, userID, isAuth := AuthenticateUser(w, r)
	if !isAuth || userID == -1 || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	llmName := parts[2]
	modelsHandler(w, r, llmName, llmKeys)
}

func UpdateLLMAPIKeys(w http.ResponseWriter, r *http.Request) {
	_, userID, isAuth := AuthenticateUser(w, r)
	if !isAuth {
		return
	}
	db, err := sql.Open("postgres", DB_URL)
	defer db.Close()
	user, err := GetUser(db, int64(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	existingLLMKeys, err := GetLLMKey(int64(userID))

	var llmkeyInput types.LLMAPIKeyInput
	json.NewDecoder(r.Body).Decode(&llmkeyInput)
	llmKeys := types.LLMAPIKey{
		Openai:      llmkeyInput.Openai,
		Palm2:       llmkeyInput.Palm2,
		Anthropic:   llmkeyInput.Anthropic,
		CohereAI:    llmkeyInput.CohereAI,
		HuggingChat: llmkeyInput.HuggingChat,
	}
	updatedLLMKeys, err := UpdateLLMKeys(&llmKeys, existingLLMKeys, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(updatedLLMKeys)
}

func geneHandler(request types.Request, llm_name string, llmKeys map[string]string) (string, error) {
	if request.Prompt == "" {
		return "", errors.New("Prompt is required")
	}
	switch llm_name {
	case "openai":
		if llmKeys["openai"] == "" {
			return "", errors.New("openai key is required")
		}
		response := openaiGenerate(&request, llmKeys["openai"])
		return response.Response, nil
	case "palm2":
		if llmKeys["palm2"] == "" {
			return "", errors.New("palm2 key is required")
		}
		response := palm2Generate(&request, llmKeys["palm2"])
		return response.Response, nil
	case "cohereai":
		if llmKeys["cohereai"] == "" {
			return "", errors.New("cohereai key is required")
		}
		response := cohereAIGenerate(&request, llmKeys["cohereai"])
		return response, nil
	case "anthropic":
		if llmKeys["anthropic"] == "" {
			return "", errors.New("anthropic key is required")
		}
		response := anthropicGenerate(&request, llmKeys["anthropic"])
		return response, nil
	}
	return "", errors.New("Invalid llm name")
}

func modelsHandler(w http.ResponseWriter, r *http.Request, llm_name string, llmKeys map[string]string) {
	request := &types.Request{}
	json.NewDecoder(r.Body).Decode(&request)
	log.Println(request)
	if request.Prompt == "" {
		json.NewEncoder(w).Encode(&types.Response{
			Response: "Prompt is required",
		})
		return
	}
	req := &types.Request{
		Prompt:      request.Prompt,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
	}
	w.Header().Set("Content-Type", "application/json")
	switch llm_name {
	case "openai":
		if llmKeys["openai"] == "" {
			json.NewEncoder(w).Encode(&types.Response{
				Response: "openai key is required",
			})
			return
		}
		response := openaiGenerate(req, llmKeys["openai"])
		json.NewEncoder(w).Encode(&types.Response{
			Response: response.Response,
		})
	case "palm2":
		if llmKeys["palm2"] == "" {
			json.NewEncoder(w).Encode(&types.Response{
				Response: "palm2 key is required",
			})
			return
		}
		response := palm2Generate(req, llmKeys["palm2"])
		json.NewEncoder(w).Encode(&types.Response{
			Response: response.Response,
		})
	case "cohereai":
		if llmKeys["cohereai"] == "" {
			json.NewEncoder(w).Encode(&types.Response{
				Response: "cohereai key is required",
			})
			return
		}
		response := cohereAIGenerate(req, llmKeys["cohereai"])
		json.NewEncoder(w).Encode(&types.Response{
			Response: response,
		})
	}
}

func openaiGenerate(request *types.Request, apiKey string) *types.Response {
	openAIRequest := &openai.CompletionRequest{
		Prompt: request.Prompt,
	}
	if request.Model != "" {
		openAIRequest.Model = request.Model
	}
	if request.MaxTokens != 0 {
		openAIRequest.MaxTokens = request.MaxTokens
	}
	if request.Temperature != 0 {
		openAIRequest.Temperature = float32(request.Temperature)
	}
	response := llms.GenerateOpenAI(*openAIRequest, apiKey)
	return &types.Response{
		Response: response,
	}
}

func palm2Generate(request *types.Request, apiKey string) *types.Response {
	palmRequest := &palm.PromptConfig{
		Prompt: palm.TextPrompt{
			Text: request.Prompt,
		},
		MaxOutputTokens: request.MaxTokens,
		Temperature:     request.Temperature,
	}
	response := llms.GeneratePaLM2(*palmRequest, apiKey)
	return &types.Response{
		Response: response,
	}
}

func cohereAIGenerate(request *types.Request, apiKey string) string {
	maxTokens := uint(request.MaxTokens)
	temperature := float64(request.Temperature)
	params := cohere.GenerateOptions{
		Prompt:      request.Prompt,
		MaxTokens:   &maxTokens,
		Temperature: &temperature,
	}
	response := llms.GenerateCohereAI(params, apiKey)
	return response
}

func anthropicGenerate(request *types.Request, apiKey string) string {
	params := anthropic.Opts{
		Sender: anthropic.Sender{
			Prompt:   request.Prompt,
			MaxToken: request.MaxTokens,
		},
	}
	response := llms.GenerateAnthropicAI(params, apiKey)
	return response
}
