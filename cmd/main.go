package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	migrate "github.com/rubenv/sql-migrate"

	"TestKafka/internal/kafka"
	"TestKafka/internal/logg"
	"TestKafka/internal/server"
	"TestKafka/internal/storage"
)

const (
	Broker = "kafka_1:9092"
)

func main() {
	kafka.Init()
	lg := logg.Logger()
	sql, err := storage.New()
	if err != nil {
		log.Fatal("error connecting to the database", "err", err)
	}
	defer func() {
		err = sql.Close()
		if err != nil {
			lg.Error("error closing the database connection", "err", err)
		}
	}()

	err = sql.Migrate(migrate.Up)
	if err != nil {
		lg.Error("migration error", "err", err)
	}

	producer, err := kafka.NewProducer([]string{Broker})
	if err != nil {
		lg.Error("error creating the producer", "err", err)
	}
	srv := server.NewServer(lg, ":8080", sql, producer)

	go func() {
		err = kafka.Consumer([]string{Broker}, context.Background())
		if err != nil {
			lg.Error("consumer error", "err", err)
		}
	}()

	go func() {
		err = srv.Run()
		if err != nil {
			log.Fatal("Error when starting the server:", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-ch
		err = srv.Shutdown()
		if err != nil {
			lg.Error("error shutting down server gracefully")
			return
		}
		done <- true
	}()
	<-done
}
