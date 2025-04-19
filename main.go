package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	h, err := NewHandler()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/task/{id:[0-9]+}/", h.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/task/{id:[0-9]+}/", h.GetTaskHandler).Methods("GET")
	r.HandleFunc("/task/", h.GetAllTasksHandler).Methods("GET")
	r.HandleFunc("/task/{id:[0-9]+}/", h.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/task/{id:[0-9]+}/", h.DeleteTaskHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
