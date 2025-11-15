package tests

import (
	"database/sql"
	"testing"

	"moneybadgers-backend/pkg"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreateItem(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a temporary items table for testing
	_, err = db.Exec("CREATE TABLE items (id STRING PRIMARY KEY, name STRING, price DECIMAL(10, 2))")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	item := pkg.Item{ID: "test-id", Name: "Test Item", Price: 9.99}
	err = pkg.CreateItem(db, item)
	if err != nil {
		t.Errorf("failed to create item: %v", err)
	}

	fetchedItem, err := pkg.GetItem(db, "test-id")
	if err != nil {
		t.Errorf("failed to get item: %v", err)
	}

	// Assert item values
	if fetchedItem.ID != item.ID || fetchedItem.Name != item.Name || fetchedItem.Price != item.Price {
		t.Errorf("fetched item does not match expected values: %+v", fetchedItem)
	}
}

// Get an item
func TestGetItem(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE items (id STRING PRIMARY KEY, name STRING, price DECIMAL(10, 2))")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	item := pkg.Item{ID: "test-id", Name: "Test Item", Price: 9.99}
	err = pkg.CreateItem(db, item)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	fetchedItem, err := pkg.GetItem(db, "test-id")
	if err != nil {
		t.Errorf("failed to get item: %v", err)
	}

	if fetchedItem.ID != item.ID || fetchedItem.Name != item.Name || fetchedItem.Price != item.Price {
		t.Errorf("fetched item does not match expected values: %+v", fetchedItem)
	}
}

// Update an item
func TestUpdateItem(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE items (id STRING PRIMARY KEY, name STRING, price DECIMAL(10, 2))")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	item := pkg.Item{ID: "test-id", Name: "Test Item", Price: 9.99}
	err = pkg.CreateItem(db, item)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	updatedItem := pkg.Item{ID: "test-id", Name: "Updated Test Item", Price: 19.99}
	err = pkg.UpdateItem(db, updatedItem)
	if err != nil {
		t.Errorf("failed to update item: %v", err)
	}

	fetchedItem, err := pkg.GetItem(db, "test-id")
	if err != nil {
		t.Errorf("failed to get item: %v", err)
	}

	if fetchedItem.ID != updatedItem.ID || fetchedItem.Name != updatedItem.Name || fetchedItem.Price != updatedItem.Price {
		t.Errorf("fetched item does not match updated values: %+v", fetchedItem)
	}
}

// Delete an item
func TestDeleteItem(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE items (id STRING PRIMARY KEY, name STRING, price DECIMAL(10, 2))")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	item := pkg.Item{ID: "test-id", Name: "Test Item", Price: 9.99}
	err = pkg.CreateItem(db, item)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	err = pkg.DeleteItem(db, "test-id")
	if err != nil {
		t.Errorf("failed to delete item: %v", err)
	}

	_, err = pkg.GetItem(db, "test-id")
	if err == nil {
		t.Errorf("expected error when getting a deleted item, got none")
	}
}
