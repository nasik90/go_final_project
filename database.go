package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func checkDatabaseExistence(dbName string) bool {
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

func createDatabase(dbName string) {
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
}

func insertTask(task Task) (int64, error) {

	db, err := sql.Open("sqlite", DbName)
	if err != nil {
		return 0, err
	}

	defer db.Close()

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

func getTasks(dbName string) (tasks []Task, err error) {

	var task Task

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return tasks, err
	}

	defer db.Close()

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

func getTask(dbName string, id int) (task Task, err error) {

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return task, err
	}

	defer db.Close()

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

func updateTask(task Task) error {
	db, err := sql.Open("sqlite", DbName)
	if err != nil {
		return err
	}

	defer db.Close()

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

func deleteTask(id int) error {
	db, err := sql.Open("sqlite", DbName)
	if err != nil {
		return err
	}

	defer db.Close()

	queryText := "DELETE FROM scheduler WHERE id = :id"

	res, err := db.Exec(queryText,
		sql.Named("id", id))

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("Задача не найдена")
	}

	return err
}
