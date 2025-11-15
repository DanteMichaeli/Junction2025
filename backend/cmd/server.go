package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"moneybadgers-backend/pkg"

	_ "github.com/mattn/go-sqlite3"
)

var clients = make(map[chan string]bool) // Store clients
var activeBasketID string                // Current active basket ID (set on startup)

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
	// Setup database with schema and sample data
	db, basketUUID, err := pkg.SetupDatabase("./app.db")
	if err != nil {
		log.Fatal("Failed to setup database:", err)
	}
	defer db.Close()

	// Set the active basket ID
	activeBasketID = basketUUID

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
				// Return all items from database
				items, err := pkg.GetAllItems(db)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(items)
			} else {
				// Get specific item from database
				item, err := pkg.GetItem(db, id)
				if err != nil {
					http.Error(w, "Item not found", http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(item)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register the SSE endpoint
	http.HandleFunc("/events", sseHandler)

	// Register endpoint to get current basket ID
	http.HandleFunc("/current-basket", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"basketId": activeBasketID,
		})
	})

	// Register the reset-demo endpoint
	http.HandleFunc("/reset-demo", func(w http.ResponseWriter, r *http.Request) {
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

		// Reset the database and get new basket UUID
		newBasketUUID, err := pkg.ResetDatabase(db)
		if err != nil {
			log.Printf("Failed to reset database: %v", err)
			http.Error(w, "Failed to reset database", http.StatusInternalServerError)
			return
		}

		// Update the active basket ID
		activeBasketID = newBasketUUID
		log.Printf("🔄 Demo reset successfully! New basket UUID: %s", activeBasketID)

		// Return success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Database reset successfully",
			"basketId": activeBasketID,
		})
	})

	// Register the add-item-to-basket endpoint
	http.HandleFunc("/add-item-to-basket", func(w http.ResponseWriter, r *http.Request) {
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

		// Parse the request (only itemId needed, basketId is global)
		var request struct {
			ItemID string `json:"itemId"`
		}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Failed to parse request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Step 1: Fetch the item from database
		item, err := pkg.GetItem(db, request.ItemID)
		if err != nil {
			http.Error(w, "Item not found in database", http.StatusNotFound)
			return
		}

		// Step 2: Add item to the active basket
		err = pkg.AddItemToBasket(db, request.ItemID, activeBasketID)
		if err != nil {
			log.Printf("Failed to add item to basket: %v", err)
			http.Error(w, "Failed to add item to basket", http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Added item %s (%s) to basket %s", item.ID, item.Name, activeBasketID)

		// Step 3: Broadcast to frontend via SSE
		go broadcastNewItem(item)

		// Return the item data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(item)
	})

	log.Println("Server starting on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
