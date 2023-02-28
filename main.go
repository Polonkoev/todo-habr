package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Task struct {
	ID    int64  `json:"id,omitempty"`
	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`
	Status      bool   `json:"status"`
}

var tasks []Task

type getAllTasks struct {
	Task    *Task  `json:"task,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response getAllTasks
	if len(tasks) == 0 {

		w.WriteHeader(http.StatusNotFound)
		response.Message = fmt.Sprintf("Список задач пуст!")
		response.Error = fmt.Sprintf("404 not found")
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

type getTaskResponse struct {
	Task  *Task  `json:"task,omitempty"`
	Error string `json:"error,omitempty"`
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskIdParam := params["id"]

	var response getTaskResponse

	taskId, err := strconv.ParseInt(taskIdParam, 10, 64)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		response.Error = fmt.Sprintf("Bad request!")

		json.NewEncoder(w).Encode(response)
		return
	}

	task := searchTask(taskId)
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		response.Error = fmt.Sprintf("task for %d not found", taskId)

		// send message to client

		json.NewEncoder(w).Encode(response)
		return
	}
	response.Task = &task

	json.NewEncoder(w).Encode(response)

}

func searchTask(taskId int64) Task {
	for _, item := range tasks {
		if item.ID == taskId {
			return item
		}

	}

	return Task{}
}

type createTaskResponse struct {
	Task    *Task  `json:"task,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response createTaskResponse
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		response.Error = err.Error()
		response.Message = "Ошибка сервера!"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Декодируем данные JSON в структуру Task
	var task Task
	err = json.Unmarshal(body, &task)

	if err != nil {

		response.Error = err.Error()
		response.Message = "Неверный формат данных JSON!"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" || task.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = fmt.Sprintf("Bad request")
		response.Message = fmt.Sprintf("Заполните или добавьте обязательные поля")

		json.NewEncoder(w).Encode(response)
		return
	}
	task.ID = int64(rand.Intn(1000000)) //todo: Найти нормальный пакет для id
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
	response.Message = fmt.Sprintf("Задача с названием %v успешно создана", task.Title)
	json.NewEncoder(w).Encode(response)

}

type updateTaskResponse struct {
	Task    *Task  `json:"task,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskIdParam := params["id"]
	var response updateTaskResponse

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error = err.Error()
		response.Message = "Ошибка сервера!"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var task Task

	err = json.Unmarshal(body, &task)
	if err != nil {
		response.Error = err.Error()
		response.Message = "Неверный формат данных JSON!"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		response.Error = fmt.Sprintf("Неверный ID!")
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" || task.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = fmt.Sprintf("Bad request")
		response.Message = fmt.Sprintf("Заполните или добавьте обязательные поля")
		json.NewEncoder(w).Encode(response)
		return
	}

	itemById := searchTask(taskId) // находим задачу по
	if itemById.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Error = fmt.Sprint("Task not found!")
		json.NewEncoder(w).Encode(response)
		return
	}

	if itemById.Status == false && task.Status == false {
		task.Status = false
	}
	task.Status = true

	for index, item := range tasks {
		if item.ID == taskId {

			task.ID = item.ID

			tasks[index] = task
			break
		}
	}

	fmt.Println("TASK from client", task)

	response.Task = &task
	response.Message = fmt.Sprintf("Задача %v успешно изменена!", task.Title)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type deleteTaskById struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskIdParam := params["id"]

	var response deleteTaskById

	taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = fmt.Sprintf("Ошибка сервера!")
		json.NewEncoder(w).Encode(response)
		return
	}

	var itemTitle = ""
	var itemId int64 = 0
	for index, item := range tasks {

		if item.ID == taskId {
			itemTitle = item.Title
			itemId = item.ID
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}

	}
	if itemId == 0 {
		w.WriteHeader(http.StatusNotFound)
		response.Error = fmt.Sprintf("Task not found")
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Message = fmt.Sprintf(`Задача  '%s' c ID %d успешно удалена!`, itemTitle, taskId)

	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	tasks = append(tasks, Task{ID: 1, Title: "Лев", Description: "Выгулять льва, налить воды крокодилам", Status: false})
	tasks = append(tasks, Task{ID: 2, Title: "Слон", Description: "Помыть черепах, почесать слону за ушами", Status: false})

	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/task/{id}", getTask).Methods("GET")
	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task/{id}", updateTask).Methods("PUT")

	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	fmt.Println("Server run")
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
