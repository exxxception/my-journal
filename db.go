package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func createTables() error {
	const usersTable = `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL
		);
		`
	const postsTable = `
		CREATE TABLE IF NOT EXISTS posts (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
    		subject TEXT NOT NULL,
    		event TEXT NOT NULL,
    		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		`

	_, err := db.Exec(usersTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(postsTable)
	if err != nil {
		return err
	}

	return nil
}

func CreateInitialDB() error {
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err) // ERROR
	}

	user := User{Username: "admin", Password: "admin"}
	_, err := CreateUser(&user)
	if err != nil {
		return fmt.Errorf("failed to create administrator: %w", err) // ERROR
	}

	return nil
}

func OpenDB(dir string) error {
	var shouldCreate bool

	var err error

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create DB directory: %w", err) // ERROR
		}
		shouldCreate = true
	}

	db, err = sql.Open("sqlite3", dir+"/store.db")
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err) // ERROR
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err) // ERROR
	}

	if shouldCreate {
		CreateInitialDB()
	}

	return nil
}

func CloseDB() error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close DB file: %w", err) // ERROR
	}
	return nil
}
