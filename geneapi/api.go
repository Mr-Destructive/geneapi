package geneapi

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mr-destructive/geneapi/geneapi/llms"
	"github.com/mr-destructive/geneapi/geneapi/types"
	"github.com/mr-destructive/palm"
	"github.com/sashabaranov/go-openai"
)

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
	db, err := sql.Open("sqlite3", DB_PATH)
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
		Prompt: request.Prompt,
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
	}
	response := llms.GeneratePaLM2(*palmRequest, apiKey)
	return &types.Response{
		Response: response,
	}
}
