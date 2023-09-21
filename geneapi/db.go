package geneapi

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const DB_PATH = "geneapi/db.sqlite3"

func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
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

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}
	fmt.Println("pinging")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	path := "geneapi/" + "migrations"
	fmt.Println(path)
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
			fmt.Println(path)
			_, err = db.Exec(string(sql))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
