package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

const (
	dbNameDefault = "scheduler.db"
	tasksLimit    = 50
	DateTemplate  = "20060102"
	portByDefault = "7540"
)

var (
	db store
)

func main() {

	var err error

	r := chi.NewRouter()

	dbFilePath := getDbFileNameAndPath()

	if checkDatabaseExistence(dbFilePath) {
		db, err = openConnection(dbFilePath)
	} else {
		db, err = createDatabase(dbFilePath)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer db.closeConnection()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleAddTask)
	r.Get("/api/tasks", handleGetTasks)
	r.Get("/api/task", handleGetTask)
	r.Put("/api/task", handleUpdateTask)
	r.Post("/api/task/done", handleDoneTask)
	r.Delete("/api/task", handleDeleteTask)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = portByDefault
	}

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(dbFilePath)

}

func getDbFileNameAndPath() string {

	dbFilePath := os.Getenv("TODO_DBFILE")

	if dbFilePath == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		dbFilePath = filepath.Join(filepath.Dir(appPath), dbNameDefault)
	}

	return dbFilePath
}
