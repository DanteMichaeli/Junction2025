package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"moneybadgers-backend/pkg"

	_ "github.com/mattn/go-sqlite3"
)

var clients = make(map[chan string]bool) // Store clients
var itemsList = []pkg.Item{}              // In-memory items list

// CORS middleware
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// SSE Handler
func sseHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)
	clients[messageChan] = true

	defer func() {
		delete(clients, messageChan)
		close(messageChan)
	}()

	// Listen for messages
	for msg := range messageChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		flusher.Flush()
	}
}

// Broadcast new item creation to all clients
func broadcastNewItem(item pkg.Item) {
	msg, _ := json.Marshal(item)
	for client := range clients {
		client <- string(msg)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case "GET":
			id := r.URL.Query().Get("id")
			if id == "" {
				// Return all items from in-memory list
				json.NewEncoder(w).Encode(itemsList)
			} else {
				// Find specific item by ID
				var foundItem *pkg.Item
				for _, item := range itemsList {
					if item.ID == id {
						foundItem = &item
						break
					}
				}
				if foundItem == nil {
					http.Error(w, "Item not found", http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(foundItem)
			}

		case "POST":
			var item pkg.Item
			json.NewDecoder(r.Body).Decode(&item)
			
			// Append to in-memory list
			itemsList = append(itemsList, item)
			log.Printf("Added item to list: %+v (Total items: %d)", item, len(itemsList))
			
			w.WriteHeader(http.StatusCreated)

			// Notify via SSE
			go broadcastNewItem(item)

		case "PUT":
			var item pkg.Item
			json.NewDecoder(r.Body).Decode(&item)
			
			// Update in in-memory list
			found := false
			for i, existingItem := range itemsList {
				if existingItem.ID == item.ID {
					itemsList[i] = item
					found = true
					log.Printf("Updated item: %+v", item)
					break
				}
			}
			
			if !found {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			
			w.WriteHeader(http.StatusOK)

		case "DELETE":
			id := r.URL.Query().Get("id")
			
			// Remove from in-memory list
			found := false
			for i, item := range itemsList {
				if item.ID == id {
					itemsList = append(itemsList[:i], itemsList[i+1:]...)
					found = true
					log.Printf("Deleted item: %s (Total items: %d)", id, len(itemsList))
					break
				}
			}
			
			if !found {
				http.Error(w, "Item not found", http.StatusNotFound)
				return
			}
			
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register the SSE endpoint
	http.HandleFunc("/events", sseHandler)

	// Register the add-predefined-item endpoint
	http.HandleFunc("/add-predefined-item", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the item name from request
		var request struct {
			ItemName string `json:"itemName"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Define predefined items with their prices
		predefinedItems := map[string]struct {
			Name  string
			Price float64
		}{
			"red bull can": {
				Name:  "Red Bull Can",
				Price: 2.49,
			},
			"vitamin well plastic bottle": {
				Name:  "Vitamin Well Plastic Bottle",
				Price: 3.99,
			},
			"chips bag": {
				Name:  "Chips Bag",
				Price: 1.99,
			},
		}

		// Validate and get the predefined item
		itemDef, exists := predefinedItems[request.ItemName]
		if !exists {
			http.Error(w, "Invalid item name. Must be one of: 'red bull can', 'vitamin well plastic bottle', or 'chips bag'", http.StatusBadRequest)
			return
		}

		// Create a new item with generated ID
		newItem := pkg.Item{
			ID:    fmt.Sprintf("%s-%d", request.ItemName, time.Now().UnixNano()),
			Name:  itemDef.Name,
			Price: itemDef.Price,
		}

		// Append to in-memory list
		itemsList = append(itemsList, newItem)
		log.Printf("Added item to list: %+v (Total items: %d)", newItem, len(itemsList))

		// Broadcast to frontend via SSE
		go broadcastNewItem(newItem)

		// Return the created item
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newItem)
	})

	// Register the classify-item endpoint
	http.HandleFunc("/classify-item", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the image data from request body
		imageData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read image data", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Create Vision API client
		ctx := context.Background()
		classifier, err := pkg.NewVisionClassifier(ctx)
		if err != nil {
			log.Printf("Failed to create vision classifier: %v", err)
			http.Error(w, "Failed to initialize classifier", http.StatusInternalServerError)
			return
		}
		defer classifier.Close()

		// Classify the image
		result, err := classifier.ClassifyImage(ctx, imageData)
		if err != nil {
			log.Printf("Failed to classify image: %v", err)
			http.Error(w, "Failed to classify image", http.StatusInternalServerError)
			return
		}

		// Return the classification result
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Server starting on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
