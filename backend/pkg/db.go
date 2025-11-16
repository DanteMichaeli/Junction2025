package pkg

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SetupDatabase creates the database schema and inserts sample items
// Baskets are created by users, not pre-populated
func SetupDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables
	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Insert sample items only (no baskets)
	err = insertSampleItems(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("✅ Database setup complete!")
	return db, nil
}

// createTables creates all necessary database tables
func createTables(db *sql.DB) error {
	// Drop existing tables to ensure clean schema
	db.Exec(`DROP TABLE IF EXISTS item_basket;`)
	db.Exec(`DROP TABLE IF EXISTS items;`)
	db.Exec(`DROP TABLE IF EXISTS baskets;`)

	// Create Items table
	_, err := db.Exec(`CREATE TABLE items (
    id STRING PRIMARY KEY,
    name STRING NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category STRING NOT NULL
);`)
	if err != nil {
		return err
	}

	// Create Baskets table
	_, err = db.Exec(`CREATE TABLE baskets (
    basketID UUID PRIMARY KEY,
    ownerName TEXT NOT NULL,
    createDate DATETIME NOT NULL,
    completedAt DATETIME
);`)
	if err != nil {
		return err
	}

	// Create Item-Basket (join) table
	_, err = db.Exec(`CREATE TABLE item_basket (
    itemID STRING,
    basketID UUID,
    FOREIGN KEY(itemID) REFERENCES items(id),
    FOREIGN KEY(basketID) REFERENCES baskets(basketID)
);`)
	if err != nil {
		return err
	}

	log.Println("📋 Tables created successfully")
	return nil
}

// insertSampleItems inserts predefined sample items only
// Baskets are created by users when they start shopping
func insertSampleItems(db *sql.DB) error {
	// Insert sample items with categories
	_, err := db.Exec(`INSERT INTO items (id, name, price, category) VALUES
    ('red-bull', 'Red Bull', 2.49, 'Beverage'),
    ('vitamin-well-refresh', 'Vitamin Well Refresh', 2.79, 'Beverage'),
    ('estrella-chips', 'Estrella Maapähkinä Rinkula', 2.89, 'Snacks');
`)
	if err != nil {
		return err
	}

	log.Println("📦 Sample items inserted successfully")
	return nil
}

// CleanupDatabase drops all tables and clears the database
func CleanupDatabase(db *sql.DB) error {
	tables := []string{"item_basket", "items", "baskets"}

	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table)
		if err != nil {
			return err
		}
	}

	log.Println("🧹 Database cleaned up successfully")
	return nil
}
