package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sqlite/authenticate"
	"sqlite/controllers"
	"sqlite/handlers"
	"sqlite/models"
	"sqlite/tasking"
	"time"
)

type Env struct {
	db        *sql.DB
	tasks     *tasking.TaskMap
	blockList *authenticate.TokenBlockList
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

func (env *Env) SignUpView(w http.ResponseWriter, r *http.Request) {
	handlers.SignUpHandler(w, r, env.db)
}

func (env *Env) LoginView(w http.ResponseWriter, r *http.Request) {
	handlers.LoginHandler(w, r, env.db)
}

func (env *Env) ChangePasswordView(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Forbidden. `Authorization` Header Required", http.StatusForbidden)
		return
	}
	handlers.ChangePasswordHandler(w, r, env.db)
	env.blockList.AddToken(tokenString)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authorizationMiddleware(blockList *authenticate.TokenBlockList) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Forbidden. `Authorization` Header Required", http.StatusForbidden)
				return
			}
			blockList.RemoveExpiredTokens()
			err := authenticate.VerifyToken(tokenString, blockList)
			if err != nil {
				http.Error(w, "Forbidden. Invalid Token", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	db, err := models.NewDB("sqlite.db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)
	if err != nil {
		log.Panic(err)
	}

	env := &Env{
		db:        db,
		tasks:     tasking.NewTaskMap(),
		blockList: authenticate.NewTokenBlockList(),
	}

	err = controllers.CreateItemTable(env.db)
	if err != nil {
		log.Panic(err)
	}

	r := mux.NewRouter()
	//no auth required

	r.Path("/signup").HandlerFunc(env.SignUpView).Methods("POST")
	r.Path("/login").HandlerFunc(env.LoginView).Methods("POST")
	r.Use(loggingMiddleware)

	// auth required
	user := r.PathPrefix("/user").Subrouter()
	user.Path("/change_password").HandlerFunc(env.ChangePasswordView).Methods("PUT")
	user.Use(authorizationMiddleware(env.blockList))

	api := r.PathPrefix("/api").Subrouter()
	api.Use(authorizationMiddleware(env.blockList))
	api.HandleFunc("/items", env.GetItemsView).Methods("GET")
	api.HandleFunc("/item", env.GetItemView).Methods("GET")
	api.HandleFunc("/item", env.CreateItemView).Methods("POST")
	api.HandleFunc("/item", env.UpdateItemView).Methods("PUT")
	api.HandleFunc("/item", env.DeleteItemView).Methods("DELETE")
	api.HandleFunc("/long", env.LongRunningGet).Methods("GET")
	api.HandleFunc("/task", env.StopTask).Methods("DELETE")
	api.HandleFunc("/task", env.GetTasks).Methods("GET")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8005",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
