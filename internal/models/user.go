package models

import "gorm.io/gorm"

type Users struct {
	//Id          uint32 `gorm:"primaryKey"`
	gorm.Model
	Surname     string `gorm:"type:varchar(50);not null"`
	Name        string `gorm:"type:varchar(50);not null"`
	Middlename  string `gorm:"type:varchar(50);"`
	PhoneNumber string `gorm:"type:char(11);not null"`
	Email       string `gorm:"type:varchar(50);"`
	Sex         string `gorm:"type:char(1);not null"`
	Password    string `gorm:"type:text;not null"`
}

func NewUser(Surname string, Name string, Middlename string, PhoneNumber string, Email string, Sex string, Password string) *Users {
	return &Users{
		Surname:     Surname,
		Name:        Name,
		Middlename:  Middlename,
		PhoneNumber: PhoneNumber,
		Email:       Email,
		Sex:         Sex,
		Password:    Password}
}
