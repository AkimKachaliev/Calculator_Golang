package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type TaskStatus string

const (
	Pending   TaskStatus = "pending"
	Completed TaskStatus = "completed"
)

type Task struct {
	ID     string     `json:"task_id"`
	Status TaskStatus `json:"status"`
	Result int        `json:"result,omitempty"`
}

type TaskManager struct {
	sync.Mutex
	tasks map[string]*Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

func (tm *TaskManager) CreateTask(expression string) *Task {
	tm.Lock()
	defer tm.Unlock()

	taskID := fmt.Sprintf("%d", rand.Intn(10000)) // Generate a unique task ID
	task := &Task{
		ID:     taskID,
		Status: Pending,
	}
	tm.tasks[taskID] = task

	// Simulate computation time
	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second) // Simulating computation time
		result := evaluateExpression(expression)               // Evaluate the expression
		tm.Lock()
		defer tm.Unlock()
		task.Status = Completed
		task.Result = result
	}()

	return task
}

func (tm *TaskManager) GetTask(taskID string) *Task {
	tm.Lock()
	defer tm.Unlock()
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil
	}
	return task
}

func evaluateExpression(expression string) int {
	// Simulated computation of the arithmetic expression
	// This function should be replaced with the actual computation logic
	// For demonstration purposes, let's just return a random number
	return rand.Intn(100)
}

func main() {
	taskManager := NewTaskManager()

	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		expression := r.FormValue("expression")
		if expression == "" {
			http.Error(w, "Expression is required", http.StatusBadRequest)
			return
		}

		task := taskManager.CreateTask(expression)
		jsonResponse(w, http.StatusOK, task)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		taskID := r.FormValue("task_id")
		if taskID == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		task := taskManager.GetTask(taskID)
		if task == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		jsonResponse(w, http.StatusOK, task)
	})

	http.ListenAndServe(":8080", nil)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
