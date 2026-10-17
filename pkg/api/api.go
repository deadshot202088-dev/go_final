package api

import "net/http"

// password пароль из TODO_PASSWORD, задаётся один раз в Init
var password string

// Init регистрирует все обработчики
func Init(pass string) {
	password = pass

	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}
