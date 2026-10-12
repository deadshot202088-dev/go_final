package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go_final_project/pkg/api"
)

const (
	webDir      = "./web"
	defaultPort = 7540
)

func Run() error {
	port := defaultPort
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Сервер запущен на порту %d\n", port)
	return http.ListenAndServe(addr, nil)
}
