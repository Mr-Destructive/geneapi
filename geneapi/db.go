package geneapi

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var DB *sql.DB

var (
	DBHost     = os.Getenv("DB_HOST")
	DBUser     = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName     = os.Getenv("DB_NAME")
	DBURL      = os.Getenv("DB_URL")
	DBPort     = 5432
)

var DB_URL string = os.Getenv("DB_URL")

func InitDB() error {
	var connStr string
	if DBURL == "" {
		connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)
	} else {
		connStr = DBURL
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = runMigrations(db)
	if err != nil {
		return err
	}
	log.Println("Database connection established")

	DB = db
	return nil
}

func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	path := "geneapi/migrations"
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			path := filepath.Join(path, file.Name())
			sql, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = db.Exec(string(sql))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
