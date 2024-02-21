package services

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
	"log/slog"
	desc "oliva-back/internal/gen"
	libjwt "oliva-back/internal/lib/jwt"
	"oliva-back/internal/storage"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// TODO: Refactor code below
type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		surname,
		name,
		middlename,
		phone_number,
		email,
		sex,
		password string,
	) (uid uint32, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (desc.UserData, error)
	//GetUser(ctx context.Context, userID uint32) (desc.UserData, error)
}

func New(
	userSaver UserSaver,
	userProvider UserProvider,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
	}
}

//func (a *Auth) GetUser(ctx context.Context, idUser uint32) (desc.UserData, error) {
//	userData, err := a.usrProvider.GetUser(ctx, idUser)
//
//	if err != nil {
//		if errors.Is(err, storage.ErrUserNotFound) {
//			return desc.UserData{}, status.Error(codes.AlreadyExists, "Пользователь не найден")
//		}
//
//		return desc.UserData{}, status.Error(codes.Internal, "Ошибка при поиске пользователя")
//	}
//
//	return userData, nil
//}

func (a *Auth) RegisterNewUser(ctx context.Context, surname, name, middlename, phoneNumber, mail, sex, password string) (uint32, error) {
	const op = "Auth.RegisterNewUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return 0, errors.New("Хэш ошибка!")
	}
	id, err := a.usrSaver.SaveUser(ctx, surname, name, middlename, phoneNumber, mail, sex, string(passHash))
	if err != nil {
		a.log.Warn("Ошибка при сохранении пользователя", logger.Error)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) Login(ctx context.Context, login, password string) (token string, err error) {
	const op = "Auth.Login"
	s := storage.Storage{}
	user, err := a.usrProvider.User(ctx, login)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("Пользователь не найден", logger.Error)

			fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Ошибка при получении пользователя", logger.Error)

		fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		errors.New("Логин или пароль неверны")
	}

	lastId := s.GetLastEntryID()
	token, err = libjwt.GenerateToken(user.Email, lastId)
	if err != nil {
		errors.New("Ошибка создания токена")
	}

	return token, nil
}
