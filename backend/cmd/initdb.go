package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create Items table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items (
    id STRING PRIMARY KEY,
    name STRING NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);`)
	if err != nil {
		log.Fatal(err)
	}

	// Create Baskets table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS baskets (
    basketID UUID PRIMARY KEY,
    createDate DATE NOT NULL,
    status TEXT CHECK(status IN ('pending', 'canceled', 'paid')) NOT NULL
);`)
	if err != nil {
		log.Fatal(err)
	}

	// Create Item-Basket (join) table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS item_basket (
    itemID STRING,
    basketID UUID,
    FOREIGN KEY(itemID) REFERENCES items(id),
    FOREIGN KEY(basketID) REFERENCES baskets(basketID)
);`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert sample data
	_, err = db.Exec(`INSERT OR IGNORE INTO items (id, name, price) VALUES
    ('pepsi-max', 'Pepsi Max', 1.99),
    ('sunmaid-sour-raisins', 'Sunmaid Sour Raisins', 1.50);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO baskets (basketID, createDate, status) VALUES
    ('550e8400-e29b-41d4-a716-446655440000', '2025-11-14', 'pending');
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO item_basket (itemID, basketID) VALUES
    ('pepsi-max', '550e8400-e29b-41d4-a716-446655440000'),
    ('sunmaid-sour-raisins', '550e8400-e29b-41d4-a716-446655440000');
`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("SQLite database initialized!")
}
