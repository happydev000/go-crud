package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// MODEL
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var dataFile = "data.json"

// ---------------------------
// LOAD JSON
// ---------------------------
func loadData() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &books)
}

// ---------------------------
// SAVE JSON
// ---------------------------
func saveData() error {
	file, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, file, 0644)
}

// ---------------------------
// RE-INDEX IDs AFTER CRUD
// ---------------------------
func reindexIDs() {
	for i := range books {
		books[i].ID = i + 1
	}
	saveData()
}

// ---------------------------
// HELPERS
// ---------------------------
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ---------------------------
// HANDLERS
// ---------------------------

// GET /books
func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, books)
}

// GET /books/{id}
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, b := range books {
		if b.ID == id {
			writeJSON(w, http.StatusOK, b)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

// POST /books
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	var newBook Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	books = append(books, newBook)
	reindexIDs()

	writeJSON(w, http.StatusCreated, newBook)
}

// PUT /books/{id}
func updateBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var update Book
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	for i := range books {
		if books[i].ID == id {
			books[i].Title = update.Title
			books[i].Author = update.Author
			saveData()
			writeJSON(w, http.StatusOK, books[i])
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

// DELETE /books/{id}
func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i := range books {
		if books[i].ID == id {
			books = append(books[:i], books[i+1:]...)
			reindexIDs()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

// ---------------------------
// CORS
// ---------------------------
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ---------------------------
// MAIN
// ---------------------------
func main() {
	err := loadData()
	if err != nil {
		log.Fatal("Failed to load JSON:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBooksHandler(w, r)
		case http.MethodPost:
			createBookHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getBookHandler(w, r)
		case http.MethodPut:
			updateBookHandler(w, r)
		case http.MethodDelete:
			deleteBookHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := enableCORS(mux)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
