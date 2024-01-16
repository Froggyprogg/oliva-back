package user_v1

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"oliva-back/pkg/models"
	desc "oliva-back/pkg/user_v1"
	"oliva-back/pkg/utils"
)

type server struct {
	desc.UnimplementedUserv1Server
}

var (
	database *gorm.DB
)

func Register(gRPCServer *grpc.Server, db *gorm.DB) {
	desc.RegisterUserv1Server(gRPCServer, &server{})

	database = db
}

// Получение пользователя
func (s *server) GetUser(ctx context.Context, req *desc.GetRequestUser) (*desc.GetResponseUser, error) {
	idUser := req.GetIdUser()

	var user models.User
	database.First(&user, idUser)

	if user.ID == 0 {
		return &desc.GetResponseUser{}, errors.New("Invalid User ID!")
	}

	return &desc.GetResponseUser{
		IdUser:      uint32(user.ID),
		Surname:     user.Surname,
		Name:        user.Name,
		Middlename:  user.Middlename,
		PhoneNumber: user.Phone_number,
		Email:       user.Email,
		Sex:         user.Sex,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *desc.PostRequestUser) (*desc.PostResponseUser, error) {
	surname := req.GetSurname()
	name := req.GetName()
	middlename := req.GetMiddlename()
	phone_number := req.GetPhoneNumber()
	mail := req.GetEmail()
	sex := req.GetSex()
	if utils.ValidateEmail(mail) == false {
		return &desc.PostResponseUser{}, errors.New("Mail invalid or missing!")
	}

	var user models.User
	database.Where(&models.User{Phone_number: phone_number}).Or(&models.User{Email: mail}).First(&user)

	if utils.CheckEmpty(user.ID) {
		return &desc.PostResponseUser{}, errors.New("Номер телефона или электронная почта заняты!")
	}

	newUser := models.NewUser(surname, name, middlename, phone_number, mail, sex) //models.User{}
	database.Create(&newUser)

	return &desc.PostResponseUser{IdUser: uint32(newUser.ID)}, nil
}
