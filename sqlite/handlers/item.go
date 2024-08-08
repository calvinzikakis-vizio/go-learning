package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sqlite/controllers"
	"sqlite/models"
	"sqlite/tasking"
	"strconv"
)

func ValidateTaskId(w http.ResponseWriter, r *http.Request, tasks *tasking.TaskMap) (int, error) {
	id := r.URL.Query().Get("taskId")
	if id == "" {
		http.Error(w, "taskId parameter is required", http.StatusBadRequest)
		return 0, errors.New("taskId parameter is required")
	}
	//convert the task id to an integer
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "taskId parameter must be an integer", http.StatusBadRequest)
		return 0, errors.New("taskId parameter must be an integer")
	}
	// verify the task id does not exist in the task map
	_, exists := tasks.GetTask(idInt)
	if exists {
		http.Error(w, "Task already exists", http.StatusConflict)
		return 0, errors.New("Task already exists")
	}
	return idInt, nil
}

func GetItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tasks *tasking.TaskMap) {
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
	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	// below call will be wrapped in the start task function
	item, err := controllers.GetItem(db, idInt, status)

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

func GetItemsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tasks *tasking.TaskMap) {
	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	// get the items
	items, err := controllers.GetItems(db, status)
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

func CreateItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tasks *tasking.TaskMap) {
	var createItem models.CreateItem
	err := json.NewDecoder(r.Body).Decode(&createItem)
	if err != nil {
		http.Error(w, "Request body not valid JSON.", http.StatusBadRequest)
		return
	}
	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	item, err := controllers.CreateItem(db, createItem.Description, status)
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

func UpdateItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tasks *tasking.TaskMap) {
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
	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	// verify the item exists
	_, err = controllers.GetItem(db, idInt, status)
	if err != nil {
		http.Error(w, "Item does not exist.", http.StatusNotFound)
		return
	}

	err = controllers.UpdateItem(db, idInt, createItem.Description, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, tasks *tasking.TaskMap) {
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

	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	// verify the item exists
	_, err = controllers.GetItem(db, idInt, status)
	if err != nil {
		http.Error(w, "Item does not exist.", http.StatusNotFound)
		return
	}

	err = controllers.DeleteItem(db, idInt, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func StopTaskHandler(w http.ResponseWriter, r *http.Request, tasks *tasking.TaskMap) {
	id := r.URL.Query().Get("taskId")
	if id == "" {
		http.Error(w, "taskId parameter is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "taskId parameter must be an integer", http.StatusBadRequest)
		return
	}

	// get the task
	status, exists := tasks.GetTask(idInt)
	if !exists {
		http.Error(w, "Task does not exist.", http.StatusNotFound)
		return
	}

	// remove the task
	status.Cancel()
	defer tasks.RemoveTask(idInt)

	w.WriteHeader(http.StatusNoContent)
}

func LongRunningGetHandler(w http.ResponseWriter, r *http.Request, tasks *tasking.TaskMap) {
	taskId, err := ValidateTaskId(w, r, tasks)
	if err != nil {
		return
	}

	// create a new context with cancel from the old one
	// this will allow us to cancel the request if the task is done
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// create a new task
	tasks.AddTask(taskId, ctx, cancel)

	// defer the removal of the task
	defer tasks.RemoveTask(taskId)

	// get the status of the task
	status, _ := tasks.GetTask(taskId)

	// start the long running task
	res, err := http.NewRequestWithContext(status.Cxt, "GET", "https://fakeresponder.com/?sleep=10000", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(http.StatusOK)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request, tasks *tasking.TaskMap) {
	tasks.Lock()
	defer tasks.Unlock()

	keys := make([]int, len(tasks.Worker))

	i := 0
	for k := range tasks.Worker {
		keys[i] = k
		i++
	}

	err := json.NewEncoder(w).Encode(keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
