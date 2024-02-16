package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"oliva-back/internal/config"
)

var (
	db              *sql.DB
	err             error
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Storage struct {
	db *sql.DB
}

func NewDatabaseConnection(cfg *config.Config) *sql.DB {
	db, err = sql.Open("postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.DBName,
			cfg.Database.Port,
			cfg.Database.SSLMode,
			cfg.Database.Timezone))

	if err != nil {
		panic(err)
	}

	log.Println("Database connection successful")

	return db
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, surname string, name string, middlename string, phone_number string, email string, sex string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(surname, name, middlename, phone_number, email, sex, password) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, surname, name, middlename, phone_number, email, sex, passHash)

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetLastEntryID() uint {
	const op = "storage.postgres.SaveUser"
	var id uint

	rows, _ := db.Query("SELECT id FROM users ORDER BY id DESC")
	defer rows.Close()

	err := rows.Scan(&id)
	if err != nil {
		fmt.Println("Ошибка при сканировании строк:", err)
	}

	return id
}
