package geneapi

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
