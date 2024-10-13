package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	taskPackage "github.com/nasik90/go_final_project/internal/entities"
	"github.com/nasik90/go_final_project/internal/storage"
)

type ResponseErr struct {
	Error string `json:"error"`
}

func ApiHandlers(store *storage.Store) *chi.Mux {

	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	//r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleAddTask)
	// r.Get("/api/tasks", handleGetTasks)
	// r.Get("/api/task", handleGetTask)
	// r.Put("/apsi/task", handleUpdateTask)
	// r.Post("/api/task/done", handleDoneTask)
	// r.Delete("/api/task", handleDeleteTask)

	r.HandleFunc("/api/task", handleAddTask())

	// mux := http.ServeMux
	// mux.HandleFunc

	return r

}

func handleNextDate(res http.ResponseWriter, req *http.Request) {
	now := req.URL.Query().Get("now")
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nowDate, _ := time.Parse(taskPackage.DateTemplate, now)
	nextDate, _ := NextDate(nowDate, date, repeat)
	res.Write([]byte(nextDate))
}

func handleAddTask(res http.ResponseWriter, req *http.Request, i int) {
	var (
		task taskPackage.Task
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
		task.Date = time.Now().Format(task.DateTemplate)
	}

	if err == nil {
		err = checkAddingTask(&task)
	}

	var dateTime, nowDate time.Time
	if err == nil {
		dateTime, err = time.Parse(taskPackage.DateTemplate, task.Date)
		nowString := time.Now().Format(taskPackage.DateTemplate)
		nowDate, _ = time.Parse(taskPackage.DateTemplate, nowString)
	}

	if err == nil && dateTime.Before(nowDate) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(taskPackage.DateTemplate)
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		}
	}

	var taskId int64
	if err == nil {
		taskId, err = storage.InsertTask(DB, task)
	}

	var responseOk ResponseOk
	if err == nil {
		responseOk.Id = strconv.Itoa(int(taskId))
		resp, err = json.Marshal(responseOk)
	}

	if err != nil {
		var responseErr ResponseErr
		responseErr.Error = err.Error()
		resp, err = json.Marshal(responseErr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusCreated)
	res.Write(resp)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	err := checkNextDateArgs(date, repeat)
	if err != nil {
		return "", err
	}

	dateTime, _ := time.Parse(task.DateTemplate, date)

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

	return newDate.Format(task.DateTemplate), nil
}

// func handleGetTasks(res http.ResponseWriter, req *http.Request) {
// 	var resp []byte
// 	type ResponseTasks struct {
// 		Tasks []Task `json:"tasks"`
// 	}

// 	tasks, err := getTasks(DbName)

// 	if err == nil {
// 		var responseTasks ResponseTasks
// 		if tasks == nil {
// 			responseTasks.Tasks = make([]Task, 0)
// 			resp, err = json.Marshal(responseTasks)
// 		} else {
// 			responseTasks.Tasks = tasks
// 			resp, err = json.Marshal(responseTasks)
// 		}
// 	}

// 	if err != nil {
// 		var responseErr ResponseErr
// 		responseErr.Error = err.Error()
// 		resp, err = json.Marshal(responseErr)
// 		if err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write(resp)
// }

// func handleGetTask(res http.ResponseWriter, req *http.Request) {
// 	var (
// 		resp []byte
// 		task Task
// 		err  error
// 		id   int
// 	)

// 	idString := req.URL.Query().Get("id")
// 	if idString == "" {
// 		err = errors.New("Не указан идентификатор")
// 	}

// 	if err == nil {
// 		id, err = strconv.Atoi(idString)
// 	}

// 	if err == nil {
// 		task, err = getTask(DbName, id)
// 	}

// 	if err == nil {
// 		resp, err = json.Marshal(task)
// 	}

// 	if err != nil {
// 		var responseErr ResponseErr
// 		responseErr.Error = err.Error()
// 		resp, err = json.Marshal(responseErr)
// 		if err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write(resp)
// }

// func handleUpdateTask(res http.ResponseWriter, req *http.Request) {
// 	var (
// 		task Task
// 		buf  bytes.Buffer
// 		resp []byte
// 		err  error
// 	)
// 	type ResponseOk struct{}

// 	_, err = buf.ReadFrom(req.Body)

// 	if err == nil {
// 		err = json.Unmarshal(buf.Bytes(), &task)
// 	}

// 	if err == nil && task.Date == "" {
// 		task.Date = time.Now().Format(DateTemplate)
// 	}

// 	err = checkAddingTask(&task)

// 	dateTime, _ := time.Parse(DateTemplate, task.Date)
// 	nowString := time.Now().Format(DateTemplate)
// 	nowDate, _ := time.Parse(DateTemplate, nowString)

// 	if err == nil && dateTime.Before(nowDate) {
// 		if task.Repeat == "" {
// 			task.Date = time.Now().Format(DateTemplate)
// 		} else {
// 			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
// 		}
// 	}

// 	if err == nil {
// 		err = updateTask(task)
// 	}

// 	if err == nil {
// 		var responseOk ResponseOk
// 		resp, err = json.Marshal(responseOk)
// 	}

// 	if err != nil {
// 		var responseErr ResponseErr
// 		responseErr.Error = err.Error()
// 		resp, err = json.Marshal(responseErr)
// 		if err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write(resp)
// }

// func handleDoneTask(res http.ResponseWriter, req *http.Request) {
// 	var (
// 		resp []byte
// 		task Task
// 		err  error
// 		id   int
// 	)

// 	type ResponseOk struct{}

// 	idString := req.URL.Query().Get("id")
// 	if idString == "" {
// 		err = errors.New("Не указан идентификатор")
// 	}

// 	if err == nil {
// 		id, err = strconv.Atoi(idString)
// 	}

// 	if err == nil {
// 		task, err = getTask(DbName, id)
// 	}

// 	if err == nil {
// 		if task.Repeat == "" {
// 			err = deleteTask(id)
// 		} else {
// 			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
// 			if err == nil {
// 				updateTask(task)
// 			}
// 		}
// 	}

// 	if err == nil {
// 		var responseOk ResponseOk
// 		resp, err = json.Marshal(responseOk)
// 	}

// 	if err != nil {
// 		var responseErr ResponseErr
// 		responseErr.Error = err.Error()
// 		resp, err = json.Marshal(responseErr)
// 		if err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write(resp)
// }

// func handleDeleteTask(res http.ResponseWriter, req *http.Request) {
// 	var (
// 		resp []byte
// 		err  error
// 		id   int
// 	)

// 	type ResponseOk struct{}

// 	idString := req.URL.Query().Get("id")
// 	if idString == "" {
// 		err = errors.New("Не указан идентификатор")
// 	}

// 	if err == nil {
// 		id, err = strconv.Atoi(idString)
// 	}

// 	if err == nil {
// 		err = deleteTask(id)
// 	}

// 	if err == nil {
// 		var responseOk ResponseOk
// 		resp, err = json.Marshal(responseOk)
// 	}

// 	if err != nil {
// 		var responseErr ResponseErr
// 		responseErr.Error = err.Error()
// 		resp, err = json.Marshal(responseErr)
// 		if err != nil {
// 			http.Error(res, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	res.WriteHeader(http.StatusCreated)
// 	res.Write(resp)
// }
