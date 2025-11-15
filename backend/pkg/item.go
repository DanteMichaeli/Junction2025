package pkg

import (
	"database/sql"
	"time"
)

// Item represents a store item
type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Basket represents a shopping basket
type Basket struct {
	ID         string    `json:"id"`
	OwnerName  string    `json:"ownerName"`
	CreateDate time.Time `json:"createDate"`
	Status     string    `json:"status"`
}

// GetItem fetches item from database by id
func GetItem(db *sql.DB, id string) (Item, error) {
	var item Item
	err := db.QueryRow("SELECT id, name, price FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name, &item.Price)
	return item, err
}

// GetAllItems fetches all items from database
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

// AddItemToBasket adds item to the current basket
func AddItemToBasket(db *sql.DB, itemID string, basketID string) error {
	_, err := db.Exec("INSERT INTO item_basket (itemID, basketID) VALUES (?, ?)", itemID, basketID)
	return err
}
