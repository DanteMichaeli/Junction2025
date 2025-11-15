package pkg

import "database/sql"

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

func UpdateItem(db *sql.DB, item Item) error {
	_, err := db.Exec("UPDATE items SET name = ?, price = ? WHERE id = ?", item.Name, item.Price, item.ID)
	return err
}

func DeleteItem(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}
