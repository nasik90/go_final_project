package main

import (
	"errors"
	"fmt"
	"time"
)

func checkAddingTask(task *Task) error {

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
