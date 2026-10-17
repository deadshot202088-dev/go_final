package main

import (
	"log"
	"os"

	"go_final_project/pkg/config"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// run выполняет всю работу приложения, чтобы отложенные вызовы
// отработали до выхода из процесса
func run() error {
	cfg := config.Load()

	if err := db.Init(cfg.DBFile); err != nil {
		return err
	}
	defer db.Close()

	return server.Run(cfg)
}
