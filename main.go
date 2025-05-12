package main

import (
	"log"
	"net/http"

	"restapi/handler"

	"github.com/gorilla/mux"
)

func main() {
	h, err := handler.NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", h.LoginHandler).Methods("POST")

	api := r.NewRoute().Subrouter()
	api.Use(h.AuthorizationMiddleware)

	api.HandleFunc("/tasks/{id:[0-9]+}", h.CreateTaskHandler).Methods("POST")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.GetTaskHandler).Methods("GET")
	api.HandleFunc("/tasks", h.GetAllTasksHandler).Methods("GET")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.UpdateTaskHandler).Methods("PUT")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.DeleteTaskHandler).Methods("DELETE")

	log.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
