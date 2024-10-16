package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ResponseErr struct {
	Error string `json:"error"`
}

func handleNextDate(res http.ResponseWriter, req *http.Request) {
	var (
		resp     []byte
		nextDate string
	)
	now := req.URL.Query().Get("now")
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nowDate, err := time.Parse(DateTemplate, now)
	if err == nil {
		nextDate, err = NextDate(nowDate, date, repeat)
	}
	if err != nil {
		resp = errorResponse(res, err)
	}
	resp = []byte(nextDate)
	writeResponse(res, resp)
}

func handleAddTask(res http.ResponseWriter, req *http.Request) {
	var (
		task Task
		buf  bytes.Buffer
		resp []byte
	)
	type ResponseOk struct {
		Id string `json:"id"`
	}

	_, err := buf.ReadFrom(req.Body)
	if err == nil {
		err = json.Unmarshal(buf.Bytes(), &task)
	}

	if err == nil && task.Date == "" {
		task.Date = time.Now().Format(DateTemplate)
	}

	if err == nil {
		err = checkAddingTask(&task)
	}

	var dateTime, nowDate time.Time
	if err == nil {
		dateTime, err = time.Parse(DateTemplate, task.Date)
		nowString := time.Now().Format(DateTemplate)
		nowDate, _ = time.Parse(DateTemplate, nowString)
	}

	if err == nil && dateTime.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(DateTemplate)
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		}
	}

	var taskId int64
	if err == nil {
		taskId, err = db.insertTask(task)
	}

	var responseOk ResponseOk
	if err == nil {
		responseOk.Id = strconv.Itoa(int(taskId))
		resp, err = json.Marshal(responseOk)
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func handleGetTasks(res http.ResponseWriter, req *http.Request) {
	var resp []byte
	type ResponseTasks struct {
		Tasks []Task `json:"tasks"`
	}

	searchingText := req.URL.Query().Get("search")

	tasks, err := db.getTasks(tasksLimit, searchingText)

	if err == nil {
		var responseTasks ResponseTasks
		if tasks == nil {
			responseTasks.Tasks = make([]Task, 0)
			resp, err = json.Marshal(responseTasks)
		} else {
			responseTasks.Tasks = tasks
			resp, err = json.Marshal(responseTasks)
		}
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func handleGetTask(res http.ResponseWriter, req *http.Request) {
	var (
		resp []byte
		task Task
		err  error
		id   int
	)

	idString := req.URL.Query().Get("id")
	if idString == "" {
		err = errors.New("не указан идентификатор")
	}

	if err == nil {
		id, err = strconv.Atoi(idString)
	}

	if err == nil {
		task, err = db.getTask(id)
	}

	if err == nil {
		resp, err = json.Marshal(task)
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func handleUpdateTask(res http.ResponseWriter, req *http.Request) {
	var (
		task Task
		buf  bytes.Buffer
		resp []byte
		err  error
	)
	type ResponseOk struct{}

	_, err = buf.ReadFrom(req.Body)

	if err == nil {
		err = json.Unmarshal(buf.Bytes(), &task)
	}

	if err == nil && task.Date == "" {
		task.Date = time.Now().Format(DateTemplate)
	}

	err = checkAddingTask(&task)

	dateTime, _ := time.Parse(DateTemplate, task.Date)
	nowString := time.Now().Format(DateTemplate)
	nowDate, _ := time.Parse(DateTemplate, nowString)

	if err == nil && dateTime.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(DateTemplate)
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		}
	}

	if err == nil {
		err = db.updateTask(task)
	}

	if err == nil {
		var responseOk ResponseOk
		resp, err = json.Marshal(responseOk)
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func handleDoneTask(res http.ResponseWriter, req *http.Request) {
	var (
		resp []byte
		task Task
		err  error
		id   int
	)

	type ResponseOk struct{}

	idString := req.URL.Query().Get("id")
	if idString == "" {
		err = errors.New("Не указан идентификатор")
	}

	if err == nil {
		id, err = strconv.Atoi(idString)
	}

	if err == nil {
		task, err = db.getTask(id)
	}

	if err == nil {
		if task.Repeat == "" {
			err = db.deleteTask(id)
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err == nil {
				db.updateTask(task)
			}
		}
	}

	if err == nil {
		var responseOk ResponseOk
		resp, err = json.Marshal(responseOk)
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func handleDeleteTask(res http.ResponseWriter, req *http.Request) {
	var (
		resp []byte
		err  error
		id   int
	)

	type ResponseOk struct{}

	idString := req.URL.Query().Get("id")
	if idString == "" {
		err = errors.New("не указан идентификатор")
	}

	if err == nil {
		id, err = strconv.Atoi(idString)
	}

	if err == nil {
		err = db.deleteTask(id)
	}

	if err == nil {
		var responseOk ResponseOk
		resp, err = json.Marshal(responseOk)
	}

	if err != nil {
		resp = errorResponse(res, err)
	}

	writeResponse(res, resp)
}

func writeResponse(res http.ResponseWriter, resp []byte) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusCreated)
	_, err := res.Write(resp)
	if err != nil {
		log.Fatal(err)
	}
}

func errorResponse(res http.ResponseWriter, err error) (resp []byte) {
	var responseErr ResponseErr
	responseErr.Error = err.Error()
	resp, err = json.Marshal(responseErr)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	return resp
}
