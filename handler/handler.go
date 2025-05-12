package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"restapi/auth"
	bt "restapi/basic_types"
	cache "restapi/cache"
	db "restapi/db"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Handler struct {
	DB    db.TaskStore
	Cache cache.TaskCache
}

func NewHandler() (*Handler, error) {
	ps, err := db.NewPostgresStore()
	if err != nil {
		return nil, err
	}

	rc, err := cache.NewRedisCache()
	if err != nil {
		return nil, err
	}

	return &Handler{DB: ps, Cache: rc}, nil
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userData db.UserData

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := h.DB.CheckUser(&userData)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, db.ErrIncorrectPassword) {
			http.Error(w, "Incorrect password", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to check user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate JWT token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task bt.Task

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	task.ID = id

	if task.ID == 0 || task.Name == "" || task.Description == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.DB.AddTask(&task); err != nil {
		if errors.Is(err, db.ErrTaskAlreadyExists) {
			http.Error(w, "Task already exists", http.StatusConflict)
		} else {
			log.Printf("Failed to insert task into DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to insert task into DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.Cache.Get(id)
	if err != nil {
		log.Printf("Failed to get from cache: %v", err)
	}

	if task == nil {
		task, err = h.DB.GetTask(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				http.Error(w, fmt.Sprintf("Task %d not found", id), http.StatusNotFound)
			} else {
				log.Printf("Failed to get task from DB: %v", err)
				http.Error(w, fmt.Sprintf("Failed to get task from DB: %v", err), http.StatusInternalServerError)
			}
			return
		}

		if err = h.Cache.Set(task); err != nil {
			log.Printf("Failed to insert to cache: %v", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DB.GetAllTasks()
	if err != nil {
		log.Printf("Failed to get all tasks from DB: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get all tasks from DB: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task bt.Task

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	task.ID = id

	if task.ID == 0 || task.Name == "" || task.Description == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTask, err := h.DB.UpdateTask(&task)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task %d not found", task.ID), http.StatusNotFound)
		} else {
			log.Printf("Failed to update task in DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to update task in DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if err = h.Cache.Delete(task.ID); err != nil {
		log.Printf("Failed to delete from cache: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			http.Error(w, fmt.Sprintf("Task %d not found", id), http.StatusNotFound)
		} else {
			log.Printf("Failed to delete task from DB: %v", err)
			http.Error(w, fmt.Sprintf("Failed to delete from DB: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if err = h.Cache.Delete(id); err != nil {
		log.Printf("Failed to delete from cache: %v", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
