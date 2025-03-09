package main

import (
	"TestKafka/internal/logg"
	"TestKafka/internal/server"
	"TestKafka/internal/storage"
	"github.com/pressly/goose"
)

func main() {
	lg := logg.NewLogger()
	sql, err := storage.NewOpen()
	if err != nil {
		return
	}
	err = sql.Ping()
	if err != nil {
		return
	}
	err = goose.Up(sql, "/home/andrey/GolandProjects/TestKafka/internal/migrations")
	if err != nil {
		return
	}
	srv := server.NewServer(lg, "8080")
	err = srv.Run()
	if err != nil {
		return
	}

}
