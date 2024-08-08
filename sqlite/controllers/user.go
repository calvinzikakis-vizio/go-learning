package controllers

import (
	"context"
	"database/sql"
	"sqlite/models"
	"time"
)

func StorePasswordInfo(db *sql.DB, user models.UserDatabase) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "INSERT INTO users VALUES (?,  ?)", user.Username, user.Hash)
	return err
}

func GetPasswordInfo(db *sql.DB, username string) (models.UserDatabase, error) {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hash string
	row := db.QueryRowContext(queryCtx, "SELECT hash FROM users WHERE username = ?", username)
	err := row.Scan(&hash)
	if err != nil {
		return models.UserDatabase{}, err
	}

	return models.UserDatabase{Username: username, Hash: hash}, nil
}

func DeletePasswordInfo(db *sql.DB, username string) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "DELETE FROM users WHERE username = ?", username)
	return err
}

func ChangePasswordInfo(db *sql.DB, user models.UserDatabase) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "UPDATE users SET hash = ? WHERE username = ?", user.Hash, user.Username)
	return err
}
