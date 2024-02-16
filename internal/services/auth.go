package services

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/logger"
	"log/slog"
	"net/http"
	desc "oliva-back/internal/gen"
	libjwt "oliva-back/internal/lib/jwt"
	"oliva-back/internal/models"
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
		sex string,
		password []byte,
	) (uid uint32, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.Users, error)
	GetUser(ctx context.Context, userID uint32) (models.Users, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		tokenTTL:    tokenTTL,
	}
}
func (a *Auth) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	//idUser := req.GetIdUser()
	//
	//var user models.Users
	//database.First(&user, idUser)
	//
	//if user.Id == 0 {
	//	return &desc.GetUserResponse{}, errors.New("Неверный User ID!")
	//}
	//
	//data := &desc.UserData{
	//	Surname:     user.Surname,
	//	Name:        user.Name,
	//	Middlename:  user.Middlename,
	//	PhoneNumber: user.PhoneNumber,
	//	Email:       user.Email,
	//	Sex:         user.Sex,
	//}
	return &desc.GetUserResponse{
		Status: http.StatusOK,
		Data:   nil,
	}, nil
}

func (a *Auth) RegisterUser(ctx context.Context, surname, name, middlename, phoneNumber, mail, sex, password string) (uint32, error) {
	const op = "Auth.RegisterNewUser"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return 0, errors.New("Хэш ошибка!")
	}

	id, err := a.usrSaver.SaveUser(ctx, surname, name, middlename, phoneNumber, mail, sex, hashed)
	if err != nil {
		a.log.Warn("Ошибка при сохранении пользователя", logger.Error)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) LoginUser(ctx context.Context, login, password string) (*desc.LoginResponse, error) {
	const op = "Auth.Login"
	s := storage.Storage{}
	user, err := a.usrProvider.User(ctx, login)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("Пользователь не найден", logger.Error)

			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Ошибка при получении пользователя", logger.Error)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &desc.LoginResponse{}, errors.New("Логин или пароль неверны")
	}

	lastId := s.GetLastEntryID()
	token, err := libjwt.GenerateToken(user.Email, lastId)
	if err != nil {
		return &desc.LoginResponse{}, errors.New("Ошибка создания токена")
	}

	return &desc.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}
