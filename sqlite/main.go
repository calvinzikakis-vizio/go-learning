package main

import (
	"context"
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
	db         *sql.DB
	cancelFunc context.CancelFunc
	ctx        context.Context
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

func (env *Env) LongRequestTimeout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mocking a long request with a 5 second timeout")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	handlers.LongRequestTimeoutHandler(w, ctx)
}

func (env *Env) LongRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mocking a long request")
	handlers.LongRequestHandler(w, env.ctx)
	log.Printf("Request Ended")
}

func (env *Env) CancelRequest(w http.ResponseWriter, r *http.Request) {
	env.cancelFunc()
	w.WriteHeader(http.StatusOK)
}

func main() {
	db, err := models.NewDB("sqlite.db")
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	env := &Env{
		db:         db,
		cancelFunc: cancelFunc,
		ctx:        ctx,
	}

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

	router.HandleFunc("/request", env.LongRequest).Methods("GET")
	router.HandleFunc("/cancel", env.CancelRequest).Methods("GET")
	router.HandleFunc("/timeout", env.LongRequestTimeout).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8005",
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
