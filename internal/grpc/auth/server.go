package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-passwd/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	desc "oliva-back/internal/gen"
	libvalidate "oliva-back/internal/lib/validate"
	auth "oliva-back/internal/services"
	"oliva-back/internal/storage"
)

type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		surname,
		name,
		middlename,
		phone_number,
		email,
		sex,
		password string,
	) (idUser uint32, err error)
	GetUser(
		ctx context.Context,
		idUser uint32,
	) (
		surname,
		name,
		middlename,
		phoneNumber,
		email,
		sex string)
}

type serverAPI struct {
	desc.UnimplementedAuthServer
	auth Auth
}

var (
	database *sql.DB
)

func Register(gRPCServer *grpc.Server, db *sql.DB) {
	desc.RegisterAuthServer(gRPCServer, &serverAPI{})
	database = db
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	if in.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "Введите логин!")
	}

	if in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Введите пароль!")
	}

	token, err := s.auth.Login(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Указанный логин или пароль неверны")
		}

		return nil, status.Error(codes.Internal, "Ошибка авторизации")
	}

	return &desc.LoginResponse{Token: token}, nil
}
func (s *serverAPI) Register(
	ctx context.Context,
	in *desc.PostUserRequest,
) (*desc.PostUserResponse, error) {
	surname := in.GetSurname()
	name := in.GetName()
	middlename := in.GetMiddlename()
	phoneNumber := in.GetPhoneNumber()
	mail := in.GetEmail()
	sex := in.GetSex()
	password := in.GetPassword()

	if libvalidate.CheckEmpty(name) {
		return &desc.PostUserResponse{}, errors.New("Введите имя!")
	}
	if libvalidate.CheckEmpty(surname) {
		return &desc.PostUserResponse{}, errors.New("Введите фамилию!")
	}
	if libvalidate.ValidateEmail(mail) == false {
		return &desc.PostUserResponse{}, errors.New("Электроная почта указана неверно или отсутсвует!")
	}
	if libvalidate.ValidatePhoneNumber(phoneNumber) == false {
		return &desc.PostUserResponse{}, errors.New("Номер телефона указан неверно или отсутсвует!")
	}
	passwordValidator := validator.New(validator.MinLength(8, errors.New("Пароль слишком короткий")), validator.MaxLength(32, errors.New("Пароль слишком длинный")))
	err := passwordValidator.Validate(password)
	if err != nil {
		return &desc.PostUserResponse{}, err
	}

	uid, err := s.auth.RegisterNewUser(ctx, surname, name, middlename, phoneNumber, mail, sex, password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "Такой пользователь уже существует")
		}

		return nil, status.Error(codes.Internal, "Ошибка при регистрации пользователя")
	}

	return &desc.PostUserResponse{IdUser: uid}, nil
}
