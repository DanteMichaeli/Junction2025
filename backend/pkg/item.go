package pkg

import (
	"database/sql"
)

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func CreateItem(db *sql.DB, item Item) error {
	_, err := db.Exec("INSERT INTO items (id, name, price) VALUES (?, ?, ?)", item.ID, item.Name, item.Price)
	return err
}

func GetItem(db *sql.DB, id string) (Item, error) {
	var item Item
	err := db.QueryRow("SELECT id, name, price FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name, &item.Price)
	return item, err
}

func GetAllItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query("SELECT id, name, price FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func UpdateItem(db *sql.DB, item Item) error {
	_, err := db.Exec("UPDATE items SET name = ?, price = ? WHERE id = ?", item.Name, item.Price, item.ID)
	return err
}

func DeleteItem(db *sql.DB, id string) error {
	// Start a transaction to ensure both deletions succeed or fail together
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, remove the item from all baskets
	_, err = tx.Exec("DELETE FROM item_basket WHERE itemID = ?", id)
	if err != nil {
		return err
	}

	// Then, delete the item itself
	_, err = tx.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
