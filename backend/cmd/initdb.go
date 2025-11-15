package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Drop and recreate Items table to add keywords column
	_, err = db.Exec(`DROP TABLE IF EXISTS items;`)
	if err != nil {
		log.Fatal(err)
	}

	// Create Items table with keywords
	_, err = db.Exec(`CREATE TABLE items (
    id STRING PRIMARY KEY,
    name STRING NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    keywords TEXT NOT NULL
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

	// Insert sample data with keywords
	pepsiKeywords := strings.Join([]string{
		"pepsi", "pepsi max", "cola", "soda",
		"soft drink", "carbonated soft drinks", "soft drinks",
		"can", "aluminum can", "steel and tin cans", "tin", "cans",
		"beverage", "drink", "non-alcoholic drink", "liquid",
		"carbonated", "cylinder", "aluminum",
		"logo", "label", "black", "thirsty",
		"steel", "gadget", "plastic",
	}, ",")

	sunmaidKeywords := strings.Join([]string{
		"sun-maid", "sunmaid", "sun maid",
		"raisin", "raisins", "sour", "sour raisin",
		"golden raisins", "dried fruit",
		"snack", "snacks", "box", "fruit",
		"packaging and labeling", "label", "logo",
		"watermelon", "flavored", "natural flavors",
	}, ",")

	vitaminWellKeywords := strings.Join([]string{
		"vitamin well", "vitamin", "well", "refresh",
		"bottle", "plastic bottle", "water bottle", "glass",
		"drink", "beverage", "water", "vitamin water", "soft drink",
		"functional drink", "fluid", "liquid",
		"drinkware", "label", "bottle cap", "personal care",
		"chemical compound", "plastic",
		"b12", "c-vitamiini", "sinkki", "lemonaden", "kiivin",
		"calorie", "juoma",
	}, ",")

	estrellaKeywords := strings.Join([]string{
		"estrella",
		"maapähkinä", "rinkula", "maapähkinävoita",
		"chips", "crisps", "snack", "snack-renkait", "peanut",
		"bag", "potato chips", "salty snack",
		"ingredient", "food", "breakfast cereal", "cereal",
		"finger food", "packaging and labeling", "produce",
		"junk food", "breakfast box", "convenience food",
		"fast food", "staple food", "recipe",
		"label", "logo", "graphic design", "advertising",
		"natural foods",
		"vegan", "makean suolainen", "rouskuva", "maku",
	}, ",")

	_, err = db.Exec(`INSERT OR IGNORE INTO items (id, name, price, keywords) VALUES
    ('pepsi-max', 'Pepsi Max', 1.99, ?),
    ('sunmaid-sour-raisins', 'Sunmaid Sour Raisins', 1.50, ?),
    ('vitamin-well-refresh', 'Vitamin Well Refresh', 3.29, ?),
    ('estrella-chips', 'Estrella Maapähkinä Rinkula', 2.99, ?);
`, pepsiKeywords, sunmaidKeywords, vitaminWellKeywords, estrellaKeywords)
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
