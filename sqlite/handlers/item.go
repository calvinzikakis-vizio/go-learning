package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sqlite/controllers"
	"sqlite/models"
	"strconv"
	"time"
)

func GetItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
		return
	}

	item, err := controllers.GetItem(db, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetItemsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	items, err := controllers.GetItems(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func CreateItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var createItem models.CreateItem
	err := json.NewDecoder(r.Body).Decode(&createItem)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}

	item, err := controllers.CreateItem(db, createItem.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func UpdateItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
		return
	}

	var createItem models.CreateItem
	err = json.NewDecoder(r.Body).Decode(&createItem)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}

	// verify the item exists
	_, err = controllers.GetItem(db, idInt)
	if err != nil {
		http.Error(w, "Item does not exist.", http.StatusNotFound)
		return
	}

	err = controllers.UpdateItem(db, idInt, createItem.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
		return
	}

	// verify the item exists
	_, err = controllers.GetItem(db, idInt)
	if err != nil {
		http.Error(w, "Item does not exist.", http.StatusNotFound)
		return
	}

	err = controllers.DeleteItem(db, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func LongRequestHandler(w http.ResponseWriter, ctx context.Context) {
	for i := 0; i < 20; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Request cancelled")
			w.WriteHeader(http.StatusRequestTimeout)
			return
		default:
			time.Sleep(1 * time.Second)
			log.Printf("Request Still Running")
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Long request complete"))
}

func LongRequestTimeoutHandler(w http.ResponseWriter, ctx context.Context) {
	for i := 0; i < 20; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Request cancelled")
			w.WriteHeader(http.StatusRequestTimeout)
			return
		default:
			time.Sleep(1 * time.Second)
			log.Printf("Request at %d seconds", i)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Long request complete"))
}
