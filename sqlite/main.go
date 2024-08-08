package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sqlite/controllers"
	"sqlite/handlers"
	"sqlite/models"
	"sqlite/tasking"
	"time"
)

type Env struct {
	db    *sql.DB
	tasks *tasking.TaskMap
}

func (env *Env) GetItemsView(w http.ResponseWriter, r *http.Request) {
	handlers.GetItemsHandler(w, r, env.db, env.tasks)
}

func (env *Env) GetItemView(w http.ResponseWriter, r *http.Request) {
	handlers.GetItemHandler(w, r, env.db, env.tasks)
}

func (env *Env) CreateItemView(w http.ResponseWriter, r *http.Request) {
	handlers.CreateItemHandler(w, r, env.db, env.tasks)
}

func (env *Env) UpdateItemView(w http.ResponseWriter, r *http.Request) {
	handlers.UpdateItemHandler(w, r, env.db, env.tasks)
}

func (env *Env) DeleteItemView(w http.ResponseWriter, r *http.Request) {
	handlers.DeleteItemHandler(w, r, env.db, env.tasks)
}

func (env *Env) StopTask(w http.ResponseWriter, r *http.Request) {
	handlers.StopTaskHandler(w, r, env.tasks)
}

func (env *Env) LongRunningGet(w http.ResponseWriter, r *http.Request) {
	handlers.LongRunningGetHandler(w, r, env.tasks)
}

func (env *Env) GetTasks(w http.ResponseWriter, r *http.Request) {
	handlers.GetTasksHandler(w, r, env.tasks)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
func main() {
	db, err := models.NewDB("sqlite.db")
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	env := &Env{
		db:    db,
		tasks: tasking.NewTaskMap(),
	}

	err = controllers.CreateItemTable(env.db)
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/items", env.GetItemsView).Methods("GET")
	router.HandleFunc("/item", env.GetItemView).Methods("GET")
	router.HandleFunc("/item", env.CreateItemView).Methods("POST")
	router.HandleFunc("/item", env.UpdateItemView).Methods("PUT")
	router.HandleFunc("/item", env.DeleteItemView).Methods("DELETE")
	router.HandleFunc("/long", env.LongRunningGet).Methods("GET")
	router.HandleFunc("/task", env.StopTask).Methods("DELETE")
	router.HandleFunc("/task", env.GetTasks).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8005",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
