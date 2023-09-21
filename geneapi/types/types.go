package types

type Request struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
}

type Model struct {
	Name        string
	Description string
}

type Response struct {
	Response string `json:"response"`
}

type LLMAPIKey struct {
	ID          int64  `json:"id"`
	Openai      string `json:"openai,omitempty"`
	Palm2       string `json:"palm2,omitempty"`
	Anthropic   string `json:"anthropic,omitempty"`
	CohereAI    string `json:"cohereai,omitempty"`
	HuggingChat string `json:"huggingchat,omitempty"`
	UserID      int64  `json:"user_id"`
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	APIKey    string    `json:"apikey"`
	LLMAPIKey LLMAPIKey `json:"llmkeys,omitempty"`
}

type LLMAPIKeyInput struct {
	Openai      string `json:"openai,omitempty"`
	Palm2       string `json:"palm2,omitempty"`
	Anthropic   string `json:"anthropic,omitempty"`
	CohereAI    string `json:"cohereai,omitempty"`
	HuggingChat string `json:"huggingchat,omitempty"`
}
