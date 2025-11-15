package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"moneybadgers-backend/pkg"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Server starting on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
