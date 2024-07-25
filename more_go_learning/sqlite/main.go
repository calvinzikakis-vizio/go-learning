package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sqlite/controllers"
	"sqlite/handlers"
	"sqlite/models"
	"time"
)

type Env struct {
	db *sql.DB
}

func (env *Env) GetItemsView(w http.ResponseWriter, r *http.Request) {
	handlers.GetItemsHandler(w, r, env.db)
}

func (env *Env) GetItemView(w http.ResponseWriter, r *http.Request) {
	handlers.GetItemHandler(w, r, env.db)
}

func (env *Env) CreateItemView(w http.ResponseWriter, r *http.Request) {
	handlers.CreateItemHandler(w, r, env.db)
}

func (env *Env) UpdateItemView(w http.ResponseWriter, r *http.Request) {
	handlers.UpdateItemHandler(w, r, env.db)
}

func (env *Env) DeleteItemView(w http.ResponseWriter, r *http.Request) {
	handlers.DeleteItemHandler(w, r, env.db)
}

func main() {
	db, err := models.NewDB("sqlite.db")
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db: db}

	err = controllers.CreateItemTable(env.db)
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/items", env.GetItemsView).Methods("GET")
	router.HandleFunc("/item", env.GetItemView).Methods("GET")
	router.HandleFunc("/item", env.CreateItemView).Methods("POST")
	router.HandleFunc("/item", env.UpdateItemView).Methods("PUT")
	router.HandleFunc("/item", env.DeleteItemView).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8005",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
