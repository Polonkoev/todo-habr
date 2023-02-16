package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`

	Description string `json:"description"`
	Done        string `json: "false"`
}

// type isDone struct {
// 	isDone string `json:"description"`
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
	err := json.NewDecoder(r.Body).Decode(&task)
	resp := make(map[string]string) // создает текст ошибки
	resp["message"] = "Bad Request" //текст ошибки

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		// // todo: send error response with appripriate http status code

		fmt.Print(err)
		json.NewEncoder(w).Encode(resp)

	}
	task.ID = strconv.Itoa(rand.Intn(1000000))
	tasks = append(tasks, task)
	// todo: http status code add
	json.NewEncoder(w).Encode(task)

	// todo: validation
	// name, descriptions check for empty

}

//

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
	fmt.Println("Server run")
	r := mux.NewRouter()
	tasks = append(tasks, Task{ID: "1", Title: "Лев", Description: "Выгулять льва, налить воды крокодилам", Done: "false"})
	tasks = append(tasks, Task{ID: "2", Title: "Слон", Description: "Помыть черепах, почесать слону за ушами", Done: "false"})
	//TODO:PATCH добавить
	// r.HandleFunc("/", getTasks).Methods("GET")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/task/{id}", getTask).Methods("GET")
	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task/{id}", updateTask).Methods("PUT")

	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:8000", r))

}
