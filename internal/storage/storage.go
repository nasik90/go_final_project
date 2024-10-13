package storage

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	task "github.com/nasik90/go_final_project/internal/entities"
	_ "modernc.org/sqlite"
)

type (
	DbConnection *sql.DB
	Store        struct{ DbConnection *sql.DB }
)

func СheckDatabaseExistence(dbName string) bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	_, err = os.Stat(dbFile)

	databaseExists := true
	if err != nil {
		databaseExists = false
	}

	return databaseExists
}

func СreateDatabase(dbName string) (DbConnection, error) {
	_, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	queryText := "CREATE TABLE scheduler (id integer PRIMARY KEY, date, title, comment, repeat text)"

	db.Exec(queryText)

	return db, err
}

func OpenConnection(dbName string) (Store, error) {
	var s Store
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return s, err
	}
	s.DbConnection = db
	return s, err
}

func (storage *Store) CloseConnection() {
	storage.DbConnection.Close()
}

func InsertTask(db *sql.DB, task task.Task) (int64, error) {

	queryText := "INSERT INTO scheduler(date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"

	res, err := db.Exec(queryText,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getTasks(db *sql.DB, dbName string) (tasks []task.Task, err error) {

	var task task.Task

	queryText := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date"

	rows, err := db.Query(queryText)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, task)
		if err != nil {
			return tasks, err
		}
	}

	return tasks, nil

}

func getTask(db *sql.DB, dbName string, id int) (task task.Task, err error) {

	queryText := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? ORDER BY date"

	rows, err := db.Query(queryText, id)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return task, err
		}
	} else {
		err = errors.New("Задача не найдена")
		return task, err
	}

	return task, nil

}

func updateTask(db *sql.DB, task task.Task) error {

	queryText := "UPDATE scheduler SET date =:date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"

	res, err := db.Exec(queryText,
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("Задача не найдена")
	}

	return err
}

func deleteTask(db *sql.DB, id int) error {

	queryText := "DELETE FROM scheduler WHERE id = :id"

	res, err := db.Exec(queryText,
		sql.Named("id", id))

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("Задача не найдена")
	}

	return err
}
