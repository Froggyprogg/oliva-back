package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"oliva-back/internal/config"
	desc "oliva-back/internal/gen"
	"oliva-back/internal/models"
)

var (
	db              *sqlx.DB
	err             error
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Storage struct {
	db *sqlx.DB
}

func NewDatabaseConnection(cfg *config.Config) *Storage {
	db, err = sqlx.Connect("pgx",
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

	return &Storage{db: db}
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, surname string, name string, middlename string, phone_number string, email string, sex string, passHash string) (uint32, error) {
	const op = "storage.SaveUser"

	//stmt, err := db.Prepare("INSERT INTO users(surname, name, middlename, phone_number, email, sex, password) VALUES(?, ?, ?, ?, ?, ?, ?)")
	//if err != nil {
	//	return 0, fmt.Errorf("%s: %w", op, err)
	//}

	//res, err := stmt.ExecContext(ctx, surname, name, middlename, phone_number, email, sex, passHash)
	//fmt.Println(res)
	//id, err := res.LastInsertId()
	//if err != nil {
	//	return 0, fmt.Errorf("%s: %w", op, err)
	//}
	//fmt.Println(id)
	tx := s.db.MustBegin()
	tx.MustExec("INSERT INTO users (surname, name, middlename, phone_number, email, sex, password) VALUES ($1, $2, $3, $4, $5, $6, $7)", surname, name, middlename, phone_number, email, sex, passHash)

	usr := models.Users{}
	err = s.db.Get(&usr, "SELECT id FROM users ORDER BY id DESC")
	fmt.Printf("%#v\n", usr.Id)

	return usr.Id, nil
}
func (s *Storage) User(ctx context.Context, login string) (desc.UserData, error) {
	const op = "storage.User"

	stmt, err := db.Prepare("SELECT id, surname, name, middlename, phone_number, email, sex, password FROM users WHERE email = ? OR phone_number = ?")
	if err != nil {
		return desc.UserData{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, login)

	var user models.Users
	err = row.Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return desc.UserData{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return desc.UserData{}, fmt.Errorf("%s: %w", op, err)
	}

	userData := desc.UserData{
		IdUser:      user.Id,
		Surname:     user.Surname,
		Name:        user.Name,
		Middlename:  user.Middlename,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Sex:         user.Sex,
		Password:    user.Password,
	}
	return userData, nil
}

//	func (s *Storage) GetUser(ctx context.Context, userID uint32) (desc.UserData, error) {
//		const op = "storage.sqlite.User"
//
//		stmt, err := s.db.Prepare("SELECT id, surname, name, middlename, phone_number, email, sex, password FROM users WHERE email = ? OR phone_number = ?")
//		if err != nil {
//			return desc.UserData{}, fmt.Errorf("%s: %w", op, err)
//		}
//
//		row := stmt.QueryRowContext(ctx, userID)
//
//		var user models.Users
//		err = row.Scan(&user.Id, &user.Email, &user.Password)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return desc.UserData{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
//			}
//
//			return desc.UserData{}, fmt.Errorf("%s: %w", op, err)
//		}
//
//		userData := desc.UserData{
//			IdUser:      user.Id,
//			Surname:     user.Surname,
//			Name:        user.Name,
//			Middlename:  user.Middlename,
//			PhoneNumber: user.PhoneNumber,
//			Email:       user.Email,
//			Sex:         user.Sex,
//			Password:    user.Password,
//		}
//		return userData, nil
//	}

func (s *Storage) GetLastEntryID() uint32 {
	const op = "storage.GetLastEntry"

	stmt := ("SELECT id FROM users ORDER BY id DESC")

	row := db.QueryRow(stmt)

	u := &models.Users{}

	err := row.Scan(&u.Id)
	if err != nil {
		fmt.Printf("%s:%w", op, err)
	}

	return u.Id
}
