package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Task struct { //Book
	ID    string `json:"id"`
	Title string `json:"title"`
	// Description *Description `json:"description"`
	Description string `json:"description"`
}

// type Description struct {
// 	Description string `json:"description"`
// }

var tasks []Task

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range tasks {

		if item.ID == params["id"] {

			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Task{})
}
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = strconv.Itoa(rand.Intn(1000000))
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range tasks {
		if item.ID == params["id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)
			var task Task
			_ = json.NewDecoder(r.Body).Decode(&task)
			task.ID = params["id"]
			tasks = append(tasks, task)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	json.NewEncoder(w).Encode(tasks)
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range tasks {
		if item.ID == params["id"] {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	r := mux.NewRouter()
	tasks = append(tasks, Task{ID: "1", Title: "Лев", Description: "Выгулять льва, налить воды крокодилам"})
	tasks = append(tasks, Task{ID: "2", Title: "Слон", Description: "Помыть черепах, почесать слону за ушами"})

	r.HandleFunc("/", getTasks).Methods("GET")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
