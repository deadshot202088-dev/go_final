package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

const tasksLimit = 50

const searchDateFormat = "02.01.2006"

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	switch {
	case search == "":
		tasks, err = db.Tasks(tasksLimit)
	default:
		if date, parseErr := time.Parse(searchDateFormat, search); parseErr == nil {
			tasks, err = db.TasksByDate(tasksLimit, date.Format(dateFormat))
		} else {
			tasks, err = db.TasksBySearch(tasksLimit, search)
		}
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasksResponse{Tasks: tasks})
}
