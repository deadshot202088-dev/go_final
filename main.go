package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	dbFile := "scheduler.db"
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
