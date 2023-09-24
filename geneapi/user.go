package geneapi

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

    _ "github.com/lib/pq"
	"github.com/mr-destructive/geneapi/geneapi/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *types.User) (types.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	nil_user := types.User{}
	if err != nil {
		log.Fatal("failed to hash password")
		return nil_user, errors.New("failed to create user")
	}

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Printf("failed to open database: %w", err)
		return nil_user, errors.New("failed to create user")
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO user(username, email, password, api_key) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("failed to prepare insert statement %w", err)
		return nil_user, errors.New("failed to create user")
	}
	defer statement.Close()

	result, err := statement.Exec(user.Username, user.Email, hashedPassword, user.APIKey)
	if err != nil {
		if UserExists(db, user.Email, user.Username) {
			return nil_user, errors.New("user already exists")
		}
		return nil_user, errors.New("failed to create user")
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		log.Fatal("failed to get last insert ID")
		return nil_user, errors.New("failed to create user")
	}
	llmKeys, err := CreateLLMKey(&user.LLMAPIKey, user.ID)
	if err != nil {
		return nil_user, err
	}
	createdUser := types.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		APIKey:    user.APIKey,
		LLMAPIKey: llmKeys,
	}
	return *&createdUser, nil
}

func GetUser(db *sql.DB, userId int64) (*types.User, error) {
	user := types.User{}

	row := db.QueryRow("SELECT id, username, email, password FROM user WHERE id = ?", userId)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserExists(db *sql.DB, email, username string) bool {
	userByEmail, _ := UserByEmail(db, email)
	if userByEmail != nil {
		return true
	}
	userByUsername, _ := UserByUsername(db, username)
	if userByUsername != nil {
		return true
	}
	return false
}

func UserByAPIKey(db *sql.DB, apiKey string) (*types.User, error) {
	query := "SELECT id FROM user WHERE api_key = ?"
	row := db.QueryRow(query, apiKey)

	user := types.User{}
	err := row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	llmQuery := "SELECT id, openai, palm2, anthropic, user_id FROM llmapiKeys WHERE user_id = ?"
	row = db.QueryRow(llmQuery, user.ID)

	llmkeys := types.LLMAPIKey{}
	err = row.Scan(&llmkeys.ID, &llmkeys.Openai, &llmkeys.Palm2, &llmkeys.Anthropic, &llmkeys.UserID)
	user.LLMAPIKey = llmkeys

	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	return &user, nil
}

func UserByEmail(db *sql.DB, email string) (*types.User, error) {
	query := "SELECT id, email, username, password FROM user WHERE email = ?"
	row := db.QueryRow(query, email)

	user := types.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, err
	}
	return &user, nil
}

func UserByUsername(db *sql.DB, username string) (*types.User, error) {
	query := "SELECT id, username, password FROM user WHERE username = ?"
	row := db.QueryRow(query, username)

	user := types.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s not found", username)
	}
	if err != nil {
		return &user, err
	}
	return &user, nil
}
