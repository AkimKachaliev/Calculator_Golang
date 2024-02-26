package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var db *sql.DB

func SetupDatabase() error {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

func SaveExpression(expression string) (string, error) {
	stmt, err := db.Prepare("INSERT INTO expressions (expression) VALUES (?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	res, err := stmt.Exec(expression)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func StartServer() {
	http.HandleFunc("/", handleRequest)
	port := 8080
	fmt.Printf("Server started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CalculateHandler(w, r)
		return
	}

	tmpl, err := template.ParseFiles(".../html/template.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при генерации HTML", http.StatusInternalServerError)
		return
	}
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Генерация уникального идентификатора
	taskID := uuid.New().String()

	// Ваша логика обработки выражения с использованием taskID

	// Пример ответа
	responseData := struct {
		TaskID string `json:"task_id"`
	}{TaskID: taskID}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Ошибка при формировании JSON-ответа", http.StatusInternalServerError)
		return
	}
}