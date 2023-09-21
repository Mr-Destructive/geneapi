package geneapi

import (
	"database/sql"
	"errors"
	"log"
)

func CreateLLMKey(LLMKey *LLMAPIKey, userID int64) (LLMAPIKey, error) {
	nil_keys := LLMAPIKey{}
	if userID == 0 {
		return nil_keys, errors.New("failed to insert LLM key")
	}

	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Printf("failed to open database: %w", err)
		return nil_keys, errors.New("failed to insert LLM key")
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO llmapikeys(user_id, openai, palm2, anthropic, cohereai, huggingchat) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("failed to prepare insert statement %w", err)
		return nil_keys, errors.New("failed to insert LLM key")
	}
	defer statement.Close()

	result, err := statement.Exec(userID, LLMKey.Openai, LLMKey.Palm2, LLMKey.Anthropic, LLMKey.CohereAI, LLMKey.HuggingChat)
	if err != nil {
		return nil_keys, errors.New("failed to insert LLM key")
	}

	LLMKey.ID, err = result.LastInsertId()
	if err != nil {
		log.Fatal("failed to get last insert ID")
		return nil_keys, errors.New("failed to insert LLM key")
	}
	llmKeys := LLMAPIKey{
		Openai:    LLMKey.Openai,
		Palm2:     LLMKey.Palm2,
		Anthropic: LLMKey.Anthropic,
        CohereAI:  LLMKey.CohereAI,
        HuggingChat: LLMKey.HuggingChat,
	}
	return llmKeys, nil
}
