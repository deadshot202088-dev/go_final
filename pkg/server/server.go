package server

import (
	"fmt"
	"net/http"

	"go_final_project/pkg/api"
	"go_final_project/pkg/config"
)

const webDir = "./web"

func Run(cfg config.Config) error {
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init(cfg.Password)

	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Сервер запущен на порту %d\n", cfg.Port)
	return http.ListenAndServe(addr, nil)
}
