package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

const (
	DbName = "scheduler.db"
)

type Task struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func main() {

	r := chi.NewRouter()

	if !checkDatabaseExistence(DbName) {
		createDatabase(DbName)
	}

	webDir := "web"
	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleAddTask)
	r.Get("/api/tasks", handleGetTasks)
	r.Get("/api/task", handleGetTask)
	r.Put("/api/task", handleUpdateTask)
	r.Post("/api/task/done", handleDoneTask)
	r.Delete("/api/task", handleDeleteTask)

	err := http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}
}

func checkAddingTask(task *Task) error {

	var errGlobal error

	if task.Title == "" {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("title is empty"))
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("date format error"))
	}

	return errGlobal

}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	err := checkNextDateArgs(date, repeat)
	if err != nil {
		return "", err
	}

	dateTime, _ := time.Parse("20060102", date)

	m_repeat := strings.Split(repeat, " ")

	firstSymbol := m_repeat[0]

	var newDate time.Time

	if firstSymbol == "y" {
		newDate = dateTime.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
	} else if firstSymbol == "d" {
		days := m_repeat[1]
		daysInt, err := strconv.Atoi(days)
		if err != nil {
			return "", err
		}
		newDate = dateTime.AddDate(0, 0, daysInt)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysInt)
		}
	}

	return newDate.Format("20060102"), nil
}

func checkNextDateArgs(date string, repeat string) error {

	var errGlobal error

	if repeat == "" {
		errGlobal = errors.Join(errGlobal, fmt.Errorf("repeat is empty"))
	}

	_, err := time.Parse("20060102", date)
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
