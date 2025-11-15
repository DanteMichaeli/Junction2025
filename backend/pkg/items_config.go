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
			"pepsi", "pepsi max", "cola", "soda",
			"soft drink", "carbonated soft drinks", "soft drinks",
			"can", "aluminum can", "steel and tin cans", "tin", "cans",
			"beverage", "drink", "non-alcoholic drink", "liquid",
			"carbonated", "cylinder", "aluminum",
			"logo", "label", "black", "thirsty",
			"steel", "gadget", "plastic",
		},
	},
	{
		ID:    "sunmaid-sour-raisins",
		Name:  "Sunmaid Sour Raisins",
		Price: 1.50,
		Keywords: []string{
			"sun-maid", "sunmaid", "sun maid",
			"raisin", "raisins", "sour", "sour raisin",
			"golden raisins", "dried fruit",
			"snack", "snacks", "box", "fruit",
			"packaging and labeling", "label", "logo",
			"watermelon", "flavored", "natural flavors",
		},
	},
	{
		ID:    "vitamin-well-refresh",
		Name:  "Vitamin Well Refresh",
		Price: 3.29,
		Keywords: []string{
			"vitamin well", "vitamin", "well", "refresh",
			"bottle", "plastic bottle", "water bottle", "glass",
			"drink", "beverage", "water", "vitamin water", "soft drink",
			"functional drink", "fluid", "liquid",
			"drinkware", "label", "bottle cap", "personal care",
			"chemical compound", "plastic",
			"b12", "c-vitamiini", "sinkki", "lemonaden", "kiivin",
			"calorie", "juoma",
		},
	},
	{
		ID:    "estrella-chips",
		Name:  "Estrella Maapähkinä Rinkula",
		Price: 2.99,
		Keywords: []string{
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
