package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Storage struct {
	db *sql.DB
}

func NewOpen() (*sql.DB, error) {

	err := godotenv.Load("/home/andrey/GolandProjects/TestKafka/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil

}
func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreareNewMessage(cxt context.Context, message string) error {
	result, err := s.db.ExecContext(cxt, "INSERT INTO messages (message) VALUES ($1)",
		message)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows inserted")
	}
	return nil
}
