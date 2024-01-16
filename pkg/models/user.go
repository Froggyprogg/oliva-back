package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	//Id        uint32 `gorm:"primaryKey"`
	gorm.Model
	Surname      string `gorm:"type:varchar(70);not null"`
	Name         string `gorm:"type:varchar(70);not null"`
	Middlename   string `gorm:"type:varchar(70); unique; not null"`
	Phone_number string `gorm:"type:char(11); unique; not null"`
	Email        string `gorm:"type:varchar(100); unique; not null"`
	Sex          string `gorm:"type:char(1); not null"`
}

func NewUser(Surname string, Name string, Middlename string, Phone_number string, Email string, Sex string) *User {
	return &User{Surname: Surname, Name: Name, Middlename: Middlename, Phone_number: Phone_number, Email: Email, Sex: Sex}
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.ID).Delete(&User{})
	return
}
