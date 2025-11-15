package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"moneybadgers-backend/pkg"

	_ "github.com/mattn/go-sqlite3"
)

var clients = make(map[chan string]bool) // Store clients

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
				items, err := pkg.GetAllItems(db)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(items)
			} else {
				item, err := pkg.GetItem(db, id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(item)
			}

		case "POST":
			var item pkg.Item
			json.NewDecoder(r.Body).Decode(&item)
			err = pkg.CreateItem(db, item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)

			// Notify via SSE
			go broadcastNewItem(item)

		case "PUT":
			var item pkg.Item
			json.NewDecoder(r.Body).Decode(&item)
			err = pkg.UpdateItem(db, item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

		case "DELETE":
			id := r.URL.Query().Get("id")
			err = pkg.DeleteItem(db, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Register the SSE endpoint
	http.HandleFunc("/events", sseHandler)

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
