package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

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

func CalculateExpression(expression string) (string, error) {
	// Парсинг математического выражения
	parts := strings.Fields(expression)

	if len(parts) != 3 {
		return "", errors.New("Неверный формат выражения. Пример: 2 + 3")
	}

	operand1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "", err
	}

	operand2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return "", err
	}

	operator := parts[1]

	var result float64
	switch operator {
	case "+":
		result = operand1 + operand2
	case "-":
		result = operand1 - operand2
	case "*":
		result = operand1 * operand2
	case "/":
		if operand2 == 0 {
			return "", errors.New("Деление на ноль")
		}
		result = operand1 / operand2
	default:
		return "", errors.New("Неподдерживаемый оператор")
	}

	return strconv.FormatFloat(result, 'f', -1, 64), nil
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
