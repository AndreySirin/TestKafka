package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type Storage struct {
	db *sql.DB
}
type Creat interface {
	CreareNewMessage(context.Context, string) error
}

func New() (*Storage, error) {
	err := godotenv.Load(".env")
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
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки соединения с БД:", err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Migrate(bool migrate.MigrationDirection) error {
	migrations := &migrate.FileMigrationSource{
		Dir: ".",
	}
	_, err := migrate.Exec(s.db, "postgres", migrations, bool)
	if err != nil {
		return fmt.Errorf("error for migrate: %v", err)
	}
	return nil
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
