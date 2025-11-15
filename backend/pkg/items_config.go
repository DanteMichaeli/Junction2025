package pkg

// PredefinedItem represents a hardcoded item for demo purposes
type PredefinedItem struct {
	ID       string
	Name     string
	Price    float64
	Keywords []string // Keywords for Vision API matching
}

// PredefinedItems contains all hardcoded items for the demo
var PredefinedItems = []PredefinedItem{
	{
		ID:    "pepsi-max",
		Name:  "Pepsi Max",
		Price: 1.99,
		Keywords: []string{
			"pepsi", "pepsi max", "cola", "soda", "soft drink",
			"can", "beverage", "drink", "carbonated",
		},
	},
	{
		ID:    "sunmaid-sour-raisins",
		Name:  "Sunmaid Sour Raisins",
		Price: 1.50,
		Keywords: []string{
			"sunmaid", "raisin", "raisins", "sour", "dried fruit",
			"snack", "box", "fruit",
		},
	},
	{
		ID:    "vitamin-well-refresh",
		Name:  "Vitamin Well Refresh",
		Price: 3.29,
		Keywords: []string{
			"vitamin well", "vitamin", "well", "refresh",
			"bottle", "drink", "beverage", "water", "vitamin water",
			"plastic bottle", "functional drink",
		},
	},
	{
		ID:    "estrella-chips",
		Name:  "Estrella Maapähkinä Rinkula",
		Price: 2.99,
		Keywords: []string{
			"estrella", "chips", "crisps", "snack", "peanut",
			"maapähkinä", "rinkula", "bag", "potato chips",
			"salty snack", "ingredient food", "breakfast cereal",
			"finger food", "packaging and labeling", "produce", "junk food", "breakfast box", "cereal",
		},
	},
}

// GetPredefinedItemByID returns a predefined item by its ID
func GetPredefinedItemByID(id string) *PredefinedItem {
	for _, item := range PredefinedItems {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

// GetAllPredefinedItems returns all predefined items
func GetAllPredefinedItems() []PredefinedItem {
	return PredefinedItems
}
