package controllers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sqlite/models"
)

func CreateItemTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, description TEXT)")
	return err
}
func GetItem(db *sql.DB, id int) (models.Item, error) {
	var item models.Item
	row := db.QueryRow("SELECT id, description FROM items WHERE id = ?", id)
	err := row.Scan(&item.Id, &item.Description)
	return item, err
}

func GetItems(db *sql.DB) ([]models.Item, error) {
	var items []models.Item
	rows, err := db.Query("SELECT id, description FROM items")
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

func UpdateItem(db *sql.DB, id int, description string) error {
	_, err := db.Exec("UPDATE items SET description = ? WHERE id = ?", description, id)
	return err
}

func DeleteItem(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}

func CreateItem(db *sql.DB, description string) (models.Item, error) {
	var id int
	// Adds an item to the item table and returns the id of the new item
	_, err := db.Exec("INSERT INTO items VALUES (NULL, ?)", description)
	if err != nil {
		return models.Item{}, err
	}

	row := db.QueryRow("SELECT last_insert_rowid()")
	err = row.Scan(&id)

	return models.Item{Id: id, Description: description}, err
}
