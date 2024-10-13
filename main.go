package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nasik90/go_final_project/internal/api"
	task "github.com/nasik90/go_final_project/internal/entities"
	"github.com/nasik90/go_final_project/internal/storage"
)

const (
	dbName = "scheduler.db"
)

var (
	//DB storage.DbConnection
	DB *sql.DB
)

// type Task struct {
// 	Id      string `json:"id"`
// 	Date    string `json:"date"`
// 	Title   string `json:"title"`
// 	Comment string `json:"comment"`
// 	Repeat  string `json:"repeat"`
// }

func main() {

	//r := chi.NewRouter()
	var store storage.Store
	//DB, _ = storage.OpenConnection(dbName)
	if storage.СheckDatabaseExistence(dbName) {
		//DB, _ = storage.OpenConnection(dbName)
		store, _ = storage.OpenConnection(dbName)
	} else {
		DB, _ = storage.СreateDatabase(dbName)
		store, _ = storage.OpenConnection(dbName)
	}
	defer store.CloseConnection()

	// webDir := "web"
	// r.Handle("/*", http.FileServer(http.Dir(webDir)))
	// //r.Get("/api/nextdate", handleNextDate)
	// r.Post("/api/task", handleAddTask)
	// // r.Get("/api/tasks", handleGetTasks)
	// // r.Get("/api/task", handleGetTask)
	// // r.Put("/apsi/task", handleUpdateTask)
	// // r.Post("/api/task/done", handleDoneTask)
	// // r.Delete("/api/task", handleDeleteTask)

	r := api.ApiHandlers(store)

	err := http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}
}

func checkAddingTask(task *task.Task) error {

	var errGlobal error

	if task.Title == "" {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("title is empty"))
	}

	_, err := time.Parse(DateTemplate, task.Date)
	if err != nil {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("date format error"))
	}

	return errGlobal

}

func checkNextDateArgs(date string, repeat string) error {

	var errGlobal error

	if repeat == "" {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("repeat is empty"))
	}

	_, err := time.Parse(DateTemplate, date)
	if err != nil {
		errGlobal = errors.Join(errGlobal, err)
	}

	repeatSplitted := strings.Split(repeat, " ")
	firstSymbol := repeatSplitted[0]
	secondSymbol := ""
	if len(repeatSplitted) > 1 {
		secondSymbol = repeatSplitted[1]
	}

	if !strings.Contains("yd", firstSymbol) {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("repeat format error"))
	}
	if strings.Contains("d", firstSymbol) && firstSymbol == "" {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("repeat format error"))
	}
	if firstSymbol == "d" {
		secondSymbolInt, err := strconv.Atoi(secondSymbol)
		if err != nil {
			errGlobal = errors.Join(errGlobal, err)
		}
		if secondSymbolInt > 400 {
			errGlobal = errors.Join(errGlobal, fmt.Errorf("repeat format error"))
		}
	}

	return errGlobal
}
