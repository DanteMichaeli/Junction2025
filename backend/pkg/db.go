package pkg

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SetupDatabase creates the database schema and inserts sample data
// Returns the database connection and the generated basket UUID
func SetupDatabase(dbPath string) (*sql.DB, string, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, "", err
	}

	// Create tables
	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, "", err
	}

	// Insert sample data and get generated basket UUID
	basketUUID, err := insertSampleData(db)
	if err != nil {
		db.Close()
		return nil, "", err
	}

	log.Println("✅ Database setup complete!")
	log.Printf("🛒 Active Basket UUID: %s", basketUUID)
	return db, basketUUID, nil
}

// createTables creates all necessary database tables
func createTables(db *sql.DB) error {
	// Create Items table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS items (
    id STRING PRIMARY KEY,
    name STRING NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);`)
	if err != nil {
		return err
	}

	// Create Baskets table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS baskets (
    basketID UUID PRIMARY KEY,
    ownerName TEXT NOT NULL,
    createDate DATE NOT NULL,
    status TEXT CHECK(status IN ('pending', 'canceled', 'paid')) NOT NULL
);`)
	if err != nil {
		return err
	}

	// Create Item-Basket (join) table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS item_basket (
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

// insertSampleData inserts predefined sample items and baskets
// Returns the generated basket UUID
func insertSampleData(db *sql.DB) (string, error) {
	// Insert sample items
	_, err := db.Exec(`INSERT OR IGNORE INTO items (id, name, price) VALUES
    ('pepsi-max', 'Pepsi Max', 1.99),
    ('sunmaid-sour-raisins', 'Sunmaid Sour Raisins', 1.50),
    ('vitamin-well-refresh', 'Vitamin Well Refresh', 3.29),
    ('estrella-chips', 'Estrella Maapähkinä Rinkula', 2.99);
`)
	if err != nil {
		return "", err
	}

	// Generate a random UUID for the basket
	basketUUID := uuid.New().String()

	// Insert sample basket with generated UUID
	_, err = db.Exec(`INSERT OR IGNORE INTO baskets (basketID, ownerName, createDate, status) VALUES
    (?, 'Demo User', date('now'), 'pending');
`, basketUUID)
	if err != nil {
		return "", err
	}

	log.Println("📦 Sample data inserted successfully")
	return basketUUID, nil
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

// ResetDatabase cleans and recreates the database with sample data
// Returns the newly generated basket UUID
func ResetDatabase(db *sql.DB) (string, error) {
	err := CleanupDatabase(db)
	if err != nil {
		return "", err
	}

	err = createTables(db)
	if err != nil {
		return "", err
	}

	basketUUID, err := insertSampleData(db)
	if err != nil {
		return "", err
	}

	log.Println("🔄 Database reset complete!")
	return basketUUID, nil
}
