package geneapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/mr-destructive/geneapi/geneapi/types"
)

type UserInfo struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	types.LLMAPIKey
}

// Generate random API key
func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
	//takee user info
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u UserInfo
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	llmkeys := types.LLMAPIKey{
		Openai:      u.LLMAPIKey.Openai,
		Palm2:       u.LLMAPIKey.Palm2,
		Anthropic:   u.LLMAPIKey.Anthropic,
		CohereAI:    u.LLMAPIKey.CohereAI,
		HuggingChat: u.LLMAPIKey.HuggingChat,
	}
	user := types.User{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		LLMAPIKey: llmkeys,
	}
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	user.APIKey = generateAPIKey()
	//create user
	user, err = CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) (map[string]string, int, bool) {
	//takee user info
	llmKeys := make(map[string]string)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return llmKeys, -1, false
	}
	// check the api key
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		http.Error(w, "API key is required", http.StatusBadRequest)
		return llmKeys, -1, false
	}
	user, err := UserByAPIKey(DB, apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return llmKeys, -1, false
	}
	llmKeys, err = GetLLMKey(user.ID)
	if err != nil {
		return llmKeys, -1, false
	}

	return llmKeys, int(user.ID), true
}
