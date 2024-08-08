package controllers

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sqlite/models"
	"sqlite/tasking"
	"time"
)

func CreateItemTable(db *sql.DB) error {
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, description TEXT)")
	return err
}

func GetItem(db *sql.DB, id int, status tasking.Status) (models.Item, error) {
	queryCtx, cancel := context.WithTimeout(status.Cxt, 5*time.Second)
	defer cancel()

	var item models.Item
	row := db.QueryRowContext(queryCtx, "SELECT id, description FROM items WHERE id = ?", id)
	err := row.Scan(&item.Id, &item.Description)
	return item, err
}

func GetItems(db *sql.DB, status tasking.Status) ([]models.Item, error) {
	// Returns all items from the item table
	// cancels the request if the context is done
	queryCtx, cancel := context.WithTimeout(status.Cxt, 5*time.Second)
	defer cancel()

	var items []models.Item
	rows, err := db.QueryContext(queryCtx, "SELECT id, description FROM items")
	if err != nil {
		return items, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.Id, &item.Description)
		if err != nil {
			return items, err
		}
		items = append(items, item)
	}
	return items, nil
}

func UpdateItem(db *sql.DB, id int, description string, status tasking.Status) error {
	queryCtx, cancel := context.WithTimeout(status.Cxt, 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "UPDATE items SET description = ? WHERE id = ?", description, id)
	return err
}

func DeleteItem(db *sql.DB, id int, status tasking.Status) error {
	queryCtx, cancel := context.WithTimeout(status.Cxt, 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(queryCtx, "DELETE FROM items WHERE id = ?", id)
	return err
}

func CreateItem(db *sql.DB, description string, status tasking.Status) (models.Item, error) {
	queryCtx, cancel := context.WithTimeout(status.Cxt, 5*time.Second)
	defer cancel()

	var id int
	// Adds an item to the item table and returns the id of the new item
	_, err := db.ExecContext(queryCtx, "INSERT INTO items VALUES (NULL, ?)", description)
	if err != nil {
		return models.Item{}, err
	}

	row := db.QueryRowContext(queryCtx, "SELECT last_insert_rowid()")
	err = row.Scan(&id)

	return models.Item{Id: id, Description: description}, err
}
