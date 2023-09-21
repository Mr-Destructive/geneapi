package geneapi

import (
	"database/sql"
	"errors"
	"log"

	"github.com/mr-destructive/geneapi/geneapi/types"
)

func CreateLLMKey(LLMKey *types.LLMAPIKey, userID int64) (types.LLMAPIKey, error) {
	nil_keys := types.LLMAPIKey{}
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
	llmKeys := types.LLMAPIKey{
		ID:          LLMKey.ID,
		Openai:      LLMKey.Openai,
		Palm2:       LLMKey.Palm2,
		Anthropic:   LLMKey.Anthropic,
		CohereAI:    LLMKey.CohereAI,
		HuggingChat: LLMKey.HuggingChat,
	}
	return llmKeys, nil
}

func GetLLMKey(userID int64) (map[string]string, error) {
	nil_keys := types.LLMAPIKey{}
	res := make(map[string]string)
	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Printf("failed to open database: %w", err)
		return res, errors.New("failed to get LLM key")
	}
	defer db.Close()

	query := "SELECT openai, palm2, anthropic FROM llmapikeys WHERE user_id = ?"
	rows := db.QueryRow(query, userID)
	if err != nil {
		return res, errors.New("failed to get LLM key")
	}
	err = rows.Scan(&nil_keys.Openai, &nil_keys.Palm2, &nil_keys.Anthropic)
	if err != nil {
		return res, errors.New("failed to get LLM key")
	}
	res["openai"] = nil_keys.Openai
	res["palm2"] = nil_keys.Palm2
	res["anthropic"] = nil_keys.Anthropic
	res["cohereai"] = nil_keys.CohereAI
	res["huggingchat"] = nil_keys.HuggingChat
	return res, nil
}

func UpdateLLMKeys(LLMKey *types.LLMAPIKey, existingLLMKeys map[string]string, userID int64) (types.LLMAPIKey, error) {
	nil_keys := types.LLMAPIKey{}
	if userID == 0 {
		return nil_keys, errors.New("user is not authenticated")
	}
	if LLMKey.Anthropic == "" {
		LLMKey.Anthropic = existingLLMKeys["anthropic"]
	}
	if LLMKey.CohereAI == "" {
		LLMKey.CohereAI = existingLLMKeys["cohereai"]
	}
	if LLMKey.HuggingChat == "" {
		LLMKey.HuggingChat = existingLLMKeys["huggingchat"]
	}
	if LLMKey.Openai == "" {
		LLMKey.Openai = existingLLMKeys["openai"]
	}
	if LLMKey.Palm2 == "" {
		LLMKey.Palm2 = existingLLMKeys["palm2"]
	}

	db, err := sql.Open("sqlite3", DB_PATH)
	if err != nil {
		log.Printf("failed to open database: %w", err)
		return nil_keys, errors.New("failed to update LLM key")
	}
	defer db.Close()

	statement, err := db.Prepare("UPDATE llmapikeys SET openai = ?, palm2 = ?, anthropic = ?, cohereai = ?, huggingchat = ? WHERE user_id = ?")
	if err != nil {
		log.Printf("failed to prepare insert statement %w", err)
		return nil_keys, errors.New("failed to update LLM key")
	}
	defer statement.Close()

	result, err := statement.Exec(LLMKey.Openai, LLMKey.Palm2, LLMKey.Anthropic, LLMKey.CohereAI, LLMKey.HuggingChat, userID)
	if err != nil {
		return nil_keys, errors.New("failed to update LLM key")
	}

	LLMKey.ID, err = result.LastInsertId()
	if err != nil {
		log.Fatal("failed to get last insert ID")
		return nil_keys, errors.New("failed to update LLM key")
	}
	llmKeys := types.LLMAPIKey{
		Openai:      LLMKey.Openai,
		Palm2:       LLMKey.Palm2,
		Anthropic:   LLMKey.Anthropic,
		CohereAI:    LLMKey.CohereAI,
		HuggingChat: LLMKey.HuggingChat,
	}
	return llmKeys, nil
}
