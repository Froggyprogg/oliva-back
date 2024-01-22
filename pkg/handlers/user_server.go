package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-passwd/validator"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	desc "oliva-back/internal/gen"
	"oliva-back/internal/models"
	"oliva-back/internal/utils"
	"os"
	"time"
)

type server struct {
	desc.UnimplementedUserServer
}

var (
	jwtSecret           = []byte(os.Getenv("JWT_SECRET"))
	TokenExpireDuration = time.Hour * 24 * 7
	database            *gorm.DB
)

func Register(gRPCServer *grpc.Server) {
	desc.RegisterUserServer(gRPCServer, &server{})

}

// Получение пользователя
func (s *server) GetUser(ctx context.Context, req *desc.GetRequestUser) (*desc.GetResponseUser, error) {
	idUser := req.GetIdUser()

	var user models.User
	database.First(&user, idUser)

	if user.Id == 0 {
		return &desc.GetResponseUser{}, errors.New("Неверный User ID!")
	}

	return &desc.GetResponseUser{
		IdUser:      uint32(user.Id),
		Surname:     user.Surname,
		Name:        user.Name,
		Middlename:  user.Middlename,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Sex:         user.Sex,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *desc.PostRequestUser) (*desc.PostResponseUser, error) {
	surname := req.GetSurname()
	name := req.GetName()
	middlename := req.GetMiddlename()
	phoneNumber := req.GetPhoneNumber()
	mail := req.GetEmail()
	sex := req.GetSex()
	password := req.GetPassword()

	if utils.CheckEmpty(name) {
		return &desc.PostResponseUser{}, errors.New("Введите имя!")
	}
	if utils.CheckEmpty(surname) {
		return &desc.PostResponseUser{}, errors.New("Введите фамилию!")
	}
	if utils.ValidateEmail(mail) == false {
		return &desc.PostResponseUser{}, errors.New("Электроная почта указана неверно или отсутсвует!")
	}

	if utils.ValidatePhoneNumber(phoneNumber) == false {
		return &desc.PostResponseUser{}, errors.New("Номер телефона указан неверно или отсутсвует!")
	}

	var user models.User
	database.Where(&models.User{PhoneNumber: phoneNumber}).Or(&models.User{Email: mail}).First(&user)

	passwordValidator := validator.New(validator.MinLength(8, errors.New("Пароль слишком короткий")), validator.MaxLength(32, errors.New("Пароль слишком длинный")))
	err := passwordValidator.Validate(password)
	if err != nil {
		return &desc.PostResponseUser{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		return &desc.PostResponseUser{}, errors.New("Хэш ошибка!")
	}

	newUser := models.NewUser(surname, name, middlename, phoneNumber, mail, sex, string(hashed)) //models.User{}
	database.Create(&newUser)

	return &desc.PostResponseUser{IdUser: newUser.Id}, nil
}

func (s *server) LoginUser(ctx context.Context, req *desc.RequestLogin) (*desc.ResponseLogin, error) {
	password := req.GetPassword()

	var user models.User
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &desc.ResponseLogin{}, errors.New("Password incorrect")
	}

	token, err := utils.NewToken(user, TokenExpireDuration, jwtSecret)
	if err != nil {
		fmt.Println("failed to generate token")

		return &desc.ResponseLogin{}, errors.New("Ошибка создания токена")
	}

	return &desc.ResponseLogin{Token: token}, nil
}
