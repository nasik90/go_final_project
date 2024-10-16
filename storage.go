package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type (
	store struct{ DbConnection *sql.DB }
)

func checkDatabaseExistence(dbFilePath string) bool {

	_, err := os.Stat(dbFilePath)

	databaseExists := true
	if err != nil {
		databaseExists = false
	}

	return databaseExists
}

func createDatabase(dbFilePath string) (store, error) {

	var s store

	_, err := os.Create(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	queryText := "CREATE TABLE scheduler (id integer PRIMARY KEY, date, title, comment, repeat text)"

	db.Exec(queryText)

	s.DbConnection = db

	return s, err
}

func openConnection(dbFilePath string) (store, error) {
	var s store
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return s, err
	}
	s.DbConnection = db
	return s, err
}

func (s *store) closeConnection() {
	s.DbConnection.Close()
}

func (s *store) insertTask(task Task) (int64, error) {

	queryText := "INSERT INTO scheduler(date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"

	res, err := s.DbConnection.Exec(queryText,
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

func (s *store) getTasks(limit int, searchingText string) (tasks []Task, err error) {

	var (
		task                Task
		searchingDateString string
		querySearchingText  string
		rows                *sql.Rows
	)

	queryText := "SELECT id, date, title, comment, repeat FROM scheduler "

	if searchingText != "" {
		searchingDate, err := time.Parse("02.01.2006", searchingText)
		if err == nil {
			searchingDateString = searchingDate.Format(DateTemplate)
			querySearchingText = "WHERE date = ? "
		} else {
			querySearchingText = "WHERE title LIKE ? OR comment LIKE ? "
		}
	}

	queryText = queryText + querySearchingText + "ORDER BY date LIMIT 10"
	queryText = strings.Replace(queryText, "10", strconv.Itoa(limit), 1)

	if searchingDateString != "" {
		rows, err = s.DbConnection.Query(queryText, searchingDateString)
	} else if searchingText != "" {
		searchingTextPerc := "%" + searchingText + "%"
		rows, err = s.DbConnection.Query(queryText, searchingTextPerc, searchingTextPerc)
	} else {
		rows, err = s.DbConnection.Query(queryText)
	}

	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, task)
		if err != nil {
			return tasks, err
		}
	}

	if err = rows.Err(); err != nil {
		return tasks, err
	}

	err = rows.Close()
	if err != nil {
		return tasks, err
	}

	return tasks, nil

}

func (s *store) getTask(id int) (task Task, err error) {

	queryText := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? ORDER BY date"

	rows, err := s.DbConnection.Query(queryText, id)
	if err != nil {
		return task, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return task, err
		}
	} else {
		err = errors.New("задача не найдена")
		return task, err
	}

	if err = rows.Err(); err != nil {
		return task, err
	}

	err = rows.Close()
	if err != nil {
		return task, err
	}
	return task, nil

}

func (s *store) updateTask(task Task) error {

	queryText := "UPDATE scheduler SET date =:date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"

	res, err := s.DbConnection.Exec(queryText,
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("задача не найдена")
	}

	return err
}

func (s *store) deleteTask(id int) error {

	queryText := "DELETE FROM scheduler WHERE id = :id"

	res, err := s.DbConnection.Exec(queryText,
		sql.Named("id", id))

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("задача не найдена")
	}

	return err
}
